package netlink

import (
	"fmt"
	"runtime"
	"sync"

	vnl "github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// Client wraps a vnl.Handle scoped to a specific network namespace.
type Client struct {
	Handle    *vnl.Handle
	NsHandle  netns.NsHandle
	Namespace string
	mu        sync.Mutex
}

// NewClient creates a new netlink client. If namespace is empty, it operates
// in the current network namespace.
func NewClient(namespace string) (*Client, error) {
	c := &Client{
		Namespace: namespace,
	}

	if namespace != "" {
		nsHandle, err := netns.GetFromName(namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get namespace %q: %w", namespace, err)
		}
		c.NsHandle = nsHandle

		handle, err := vnl.NewHandleAt(nsHandle)
		if err != nil {
			nsHandle.Close()
			return nil, fmt.Errorf("failed to create netlink handle for namespace %q: %w", namespace, err)
		}
		c.Handle = handle
	} else {
		handle, err := vnl.NewHandle()
		if err != nil {
			return nil, fmt.Errorf("failed to create netlink handle: %w", err)
		}
		c.Handle = handle
		c.NsHandle = -1
	}

	return c, nil
}

// Close releases the netlink handle and namespace handle.
func (c *Client) Close() {
	if c.Handle != nil {
		c.Handle.Close()
	}
	if c.NsHandle >= 0 {
		c.NsHandle.Close()
	}
}

// RunInNamespace executes a function within the client's network namespace.
func (c *Client) RunInNamespace(fn func() error) error {
	if c.Namespace == "" {
		return fn()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origNs, err := netns.Get()
	if err != nil {
		return fmt.Errorf("failed to get current namespace: %w", err)
	}
	defer origNs.Close()

	if err := netns.Set(c.NsHandle); err != nil {
		return fmt.Errorf("failed to set namespace: %w", err)
	}
	defer netns.Set(origNs) //nolint:errcheck

	return fn()
}
