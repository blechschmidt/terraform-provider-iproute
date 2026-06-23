package netlink

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vishvananda/netns"
)

// ResolveNamespace turns a namespace specification into an open netns handle.
//
// Supported forms:
//
//	""                  -> the caller's current network namespace
//	"<name>"            -> a named namespace under /run/netns (ip netns add)
//	"name:<name>"       -> the same, stated explicitly
//	"pid:<pid>"         -> the network namespace of the given process
//	"path:<path>"       -> a nsfs path, e.g. /proc/<pid>/ns/net or a bind mount
//	"docker:<id|name>"  -> the network namespace of a running Docker container,
//	                       resolved through the Docker API (its init PID)
//
// The "docker:" form lets a Terraform module target a container created by the
// docker provider without any shell helper: reference the container id, e.g.
// namespace = "docker:${docker_container.gateway.id}". Resolution happens when
// the provider is configured, so the container must already be running by then
// (depend the provider configuration on the container by interpolating a value
// only known after it is created, such as its id).
func ResolveNamespace(spec string) (netns.NsHandle, error) {
	if spec == "" {
		return netns.Get()
	}

	kind, val, hasPrefix := strings.Cut(spec, ":")
	if !hasPrefix {
		// Bare value: a named namespace, for backwards compatibility.
		return netns.GetFromName(spec)
	}

	switch kind {
	case "name":
		return netns.GetFromName(val)
	case "path":
		return netns.GetFromPath(val)
	case "pid":
		pid, err := strconv.Atoi(val)
		if err != nil {
			return -1, fmt.Errorf("invalid pid %q in namespace %q: %w", val, spec, err)
		}
		return netns.GetFromPid(pid)
	case "docker":
		pid, err := dockerInitPID(val)
		if err != nil {
			return -1, fmt.Errorf("resolving docker namespace %q: %w", spec, err)
		}
		return netns.GetFromPid(pid)
	default:
		return -1, fmt.Errorf("unknown namespace kind %q (want name:, pid:, path: or docker:)", kind)
	}
}

// dockerInitPID returns the init PID of a running container, queried over the
// Docker Engine API on the local socket. Honours DOCKER_HOST when it names a
// unix socket; otherwise defaults to /var/run/docker.sock.
func dockerInitPID(idOrName string) (int, error) {
	socket := "/var/run/docker.sock"
	if h := os.Getenv("DOCKER_HOST"); strings.HasPrefix(h, "unix://") {
		socket = strings.TrimPrefix(h, "unix://")
	}

	httpc := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", socket)
			},
		},
	}

	url := "http://docker/containers/" + idOrName + "/json"
	resp, err := httpc.Get(url)
	if err != nil {
		return 0, fmt.Errorf("querying docker socket %s: %w", socket, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("docker inspect of %q returned HTTP %d", idOrName, resp.StatusCode)
	}

	var info struct {
		State struct {
			Pid     int  `json:"Pid"`
			Running bool `json:"Running"`
		} `json:"State"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, fmt.Errorf("decoding docker inspect of %q: %w", idOrName, err)
	}
	if !info.State.Running || info.State.Pid == 0 {
		return 0, fmt.Errorf("container %q is not running (pid %d)", idOrName, info.State.Pid)
	}
	return info.State.Pid, nil
}
