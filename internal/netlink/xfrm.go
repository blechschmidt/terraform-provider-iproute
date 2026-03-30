package netlink

import (
	vnl "github.com/vishvananda/netlink"
)

func (c *Client) XfrmStateAdd(state *vnl.XfrmState) error {
	return c.Handle.XfrmStateAdd(state)
}

func (c *Client) XfrmStateDel(state *vnl.XfrmState) error {
	return c.Handle.XfrmStateDel(state)
}

func (c *Client) XfrmStateUpdate(state *vnl.XfrmState) error {
	return c.Handle.XfrmStateUpdate(state)
}

func (c *Client) XfrmStateList(family int) ([]vnl.XfrmState, error) {
	return c.Handle.XfrmStateList(family)
}

func (c *Client) XfrmStateGet(state *vnl.XfrmState) (*vnl.XfrmState, error) {
	return c.Handle.XfrmStateGet(state)
}

func (c *Client) XfrmPolicyAdd(policy *vnl.XfrmPolicy) error {
	return c.Handle.XfrmPolicyAdd(policy)
}

func (c *Client) XfrmPolicyDel(policy *vnl.XfrmPolicy) error {
	return c.Handle.XfrmPolicyDel(policy)
}

func (c *Client) XfrmPolicyUpdate(policy *vnl.XfrmPolicy) error {
	return c.Handle.XfrmPolicyUpdate(policy)
}

func (c *Client) XfrmPolicyList(family int) ([]vnl.XfrmPolicy, error) {
	return c.Handle.XfrmPolicyList(family)
}

func (c *Client) XfrmPolicyGet(policy *vnl.XfrmPolicy) (*vnl.XfrmPolicy, error) {
	return c.Handle.XfrmPolicyGet(policy)
}
