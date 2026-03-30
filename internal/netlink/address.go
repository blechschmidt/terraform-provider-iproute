package netlink

import (
	vnl "github.com/vishvananda/netlink"
)

func (c *Client) AddrAdd(link vnl.Link, addr *vnl.Addr) error {
	return c.Handle.AddrAdd(link, addr)
}

func (c *Client) AddrDel(link vnl.Link, addr *vnl.Addr) error {
	return c.Handle.AddrDel(link, addr)
}

func (c *Client) AddrReplace(link vnl.Link, addr *vnl.Addr) error {
	return c.Handle.AddrReplace(link, addr)
}

func (c *Client) AddrList(link vnl.Link, family int) ([]vnl.Addr, error) {
	return c.Handle.AddrList(link, family)
}
