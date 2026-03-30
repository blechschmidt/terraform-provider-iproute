package netlink

import (
	vnl "github.com/vishvananda/netlink"
)

func (c *Client) NeighAdd(neigh *vnl.Neigh) error {
	return c.Handle.NeighAdd(neigh)
}

func (c *Client) NeighDel(neigh *vnl.Neigh) error {
	return c.Handle.NeighDel(neigh)
}

func (c *Client) NeighSet(neigh *vnl.Neigh) error {
	return c.Handle.NeighSet(neigh)
}

func (c *Client) NeighList(linkIndex int, family int) ([]vnl.Neigh, error) {
	return c.Handle.NeighList(linkIndex, family)
}

func (c *Client) NeighProxyList(linkIndex int, family int) ([]vnl.Neigh, error) {
	return c.Handle.NeighProxyList(linkIndex, family)
}
