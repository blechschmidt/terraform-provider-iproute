package netlink

// Segment routing management via ip command.

func (c *Client) SrTunSrcSet(address string) error {
	return c.ipExec("sr", "tunsrc", "set", address)
}

func (c *Client) SrTunSrcGet() (string, error) {
	return c.ipExecOutput("sr", "tunsrc", "show")
}
