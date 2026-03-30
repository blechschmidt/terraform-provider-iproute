package netlink

import (
	vnl "github.com/vishvananda/netlink"
)

func (c *Client) RuleAdd(rule *vnl.Rule) error {
	return c.Handle.RuleAdd(rule)
}

func (c *Client) RuleDel(rule *vnl.Rule) error {
	return c.Handle.RuleDel(rule)
}

func (c *Client) RuleList(family int) ([]vnl.Rule, error) {
	return c.Handle.RuleList(family)
}

func (c *Client) RuleListFiltered(family int, filter *vnl.Rule, mask uint64) ([]vnl.Rule, error) {
	return c.Handle.RuleListFiltered(family, filter, mask)
}
