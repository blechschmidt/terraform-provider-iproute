package netlink

// Multicast address management via ip command.

func (c *Client) MaddressAdd(device, address string) error {
	return c.ipExec("maddress", "add", address, "dev", device)
}

func (c *Client) MaddressDel(device, address string) error {
	return c.ipExec("maddress", "del", address, "dev", device)
}

func (c *Client) MaddressList(device string) (string, error) {
	if device != "" {
		return c.ipExecOutput("maddress", "show", "dev", device)
	}
	return c.ipExecOutput("maddress", "show")
}
