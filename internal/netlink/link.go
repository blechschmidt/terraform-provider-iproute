package netlink

import (
	"net"

	vnl "github.com/vishvananda/netlink"
)

func (c *Client) LinkAdd(link vnl.Link) error {
	// TUN/TAP devices are created via /dev/net/tun ioctl, not netlink messages.
	// The Handle's namespace scoping only affects netlink sockets, so we must
	// switch the OS thread to the target namespace for tuntap creation.
	if _, ok := link.(*vnl.Tuntap); ok {
		return c.RunInNamespace(func() error {
			return c.Handle.LinkAdd(link)
		})
	}
	return c.Handle.LinkAdd(link)
}

func (c *Client) LinkDel(link vnl.Link) error {
	return c.Handle.LinkDel(link)
}

func (c *Client) LinkSetUp(link vnl.Link) error {
	return c.Handle.LinkSetUp(link)
}

func (c *Client) LinkSetDown(link vnl.Link) error {
	return c.Handle.LinkSetDown(link)
}

func (c *Client) LinkSetName(link vnl.Link, name string) error {
	return c.Handle.LinkSetName(link, name)
}

func (c *Client) LinkSetMTU(link vnl.Link, mtu int) error {
	return c.Handle.LinkSetMTU(link, mtu)
}

func (c *Client) LinkSetMaster(link vnl.Link, master vnl.Link) error {
	return c.Handle.LinkSetMaster(link, master)
}

func (c *Client) LinkSetNoMaster(link vnl.Link) error {
	return c.Handle.LinkSetNoMaster(link)
}

func (c *Client) LinkSetHardwareAddr(link vnl.Link, addr net.HardwareAddr) error {
	return c.Handle.LinkSetHardwareAddr(link, addr)
}

func (c *Client) LinkSetTxQLen(link vnl.Link, qlen int) error {
	return c.Handle.LinkSetTxQLen(link, qlen)
}

func (c *Client) LinkSetAlias(link vnl.Link, alias string) error {
	return c.Handle.LinkSetAlias(link, alias)
}

func (c *Client) LinkByName(name string) (vnl.Link, error) {
	return c.Handle.LinkByName(name)
}

func (c *Client) LinkByIndex(index int) (vnl.Link, error) {
	return c.Handle.LinkByIndex(index)
}

func (c *Client) LinkList() ([]vnl.Link, error) {
	return c.Handle.LinkList()
}
