package netlink

// TCP metrics management via ip command.

func (c *Client) TcpMetricsAdd(address string) error {
	return c.ipExec("tcp_metrics", "add", address)
}

func (c *Client) TcpMetricsDel(address string) error {
	return c.ipExec("tcp_metrics", "delete", address)
}

func (c *Client) TcpMetricsList() (string, error) {
	return c.ipExecOutput("tcp_metrics", "show")
}
