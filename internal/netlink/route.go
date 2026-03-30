package netlink

import (
	"net"

	vnl "github.com/vishvananda/netlink"
)

func (c *Client) RouteAdd(route *vnl.Route) error {
	return c.Handle.RouteAdd(route)
}

func (c *Client) RouteDel(route *vnl.Route) error {
	return c.Handle.RouteDel(route)
}

func (c *Client) RouteReplace(route *vnl.Route) error {
	return c.Handle.RouteReplace(route)
}

func (c *Client) RouteList(link vnl.Link, family int) ([]vnl.Route, error) {
	return c.Handle.RouteList(link, family)
}

func (c *Client) RouteListFiltered(family int, filter *vnl.Route, mask uint64) ([]vnl.Route, error) {
	return c.Handle.RouteListFiltered(family, filter, mask)
}

func (c *Client) RouteGet(dst net.IP) ([]vnl.Route, error) {
	return c.Handle.RouteGet(dst)
}
