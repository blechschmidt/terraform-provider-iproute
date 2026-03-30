package netlink

import (
	"strconv"
	"strings"

	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// FOU operations: Add/Del use netlink via RunInNamespace (to ensure the
// generic netlink socket is in the right namespace). List uses the ip command
// because vishvananda/netlink's FouList has a parsing bug.

func (c *Client) FouAdd(f vnl.Fou) error {
	return c.RunInNamespace(func() error {
		return vnl.FouAdd(f)
	})
}

func (c *Client) FouDel(f vnl.Fou) error {
	return c.RunInNamespace(func() error {
		return vnl.FouDel(f)
	})
}

func (c *Client) FouList(family int) ([]vnl.Fou, error) {
	args := []string{"fou", "show"}
	if family == unix.AF_INET6 {
		args = append([]string{"-6"}, args...)
	}
	out, err := c.ipExecOutput(args...)
	if err != nil {
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	var fous []vnl.Fou
	for _, line := range strings.Split(out, "\n") {
		f := parseFouLine(line)
		if f.Port > 0 {
			fous = append(fous, f)
		}
	}
	return fous, nil
}

// parseFouLine parses a line like "port 5555 ipproto 4" or "port 5556 gue"
func parseFouLine(line string) vnl.Fou {
	f := vnl.Fou{Family: unix.AF_INET}
	fields := strings.Fields(line)
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "port":
			if i+1 < len(fields) {
				v, _ := strconv.Atoi(fields[i+1])
				f.Port = v
				i++
			}
		case "ipproto":
			if i+1 < len(fields) {
				v, _ := strconv.Atoi(fields[i+1])
				f.Protocol = v
				i++
			}
		case "gue":
			f.EncapType = vnl.FOU_ENCAP_GUE
		}
	}
	return f
}
