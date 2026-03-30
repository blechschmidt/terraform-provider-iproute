package netlink

import "strings"

// IPv6 token management via ip command (IFLA_INET6_TOKEN).

func (c *Client) TokenSet(device, token string) error {
	return c.ipExec("token", "set", token, "dev", device)
}

func (c *Client) TokenGet(device string) (string, error) {
	out, err := c.ipExecOutput("token", "get", "dev", device)
	if err != nil {
		return "", err
	}
	// Output format: "token ::1 dev veth0"
	out = strings.TrimSpace(out)
	parts := strings.Fields(out)
	if len(parts) >= 2 && parts[0] == "token" {
		return parts[1], nil
	}
	return out, nil
}
