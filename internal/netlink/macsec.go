package netlink

import "fmt"

// MACsec management via ip command (genetlink not directly supported by vishvananda/netlink).

type Macsec struct {
	Parent      string
	Name        string
	SCI         string
	Port        int
	Encrypt     bool
	CipherSuite string
	ICVLen      int
	EncodingSA  int
	Validate    string
	Protect     bool
	ReplayProtect bool
	Window      int
}

func (c *Client) MacsecAdd(m *Macsec) error {
	args := []string{"link", "add", "link", m.Parent, m.Name, "type", "macsec"}
	if m.Port > 0 {
		args = append(args, "port", fmt.Sprintf("%d", m.Port))
	}
	if m.SCI != "" {
		args = append(args, "sci", m.SCI)
	}
	if m.Encrypt {
		args = append(args, "encrypt", "on")
	}
	if m.CipherSuite != "" {
		args = append(args, "cipher", m.CipherSuite)
	}
	if m.ICVLen > 0 {
		args = append(args, "icvlen", fmt.Sprintf("%d", m.ICVLen))
	}
	if m.Validate != "" {
		args = append(args, "validate", m.Validate)
	}
	if m.Protect {
		args = append(args, "protect", "on")
	}
	if m.ReplayProtect {
		args = append(args, "replay", "on", "window", fmt.Sprintf("%d", m.Window))
	}
	return c.ipExec(args...)
}

func (c *Client) MacsecDel(name string) error {
	return c.ipExec("link", "del", name)
}
