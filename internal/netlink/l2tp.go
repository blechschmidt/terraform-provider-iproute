package netlink

import (
	"fmt"
)

// L2TP tunnel and session management via raw genetlink.

type L2tpTunnel struct {
	TunnelID     int
	PeerTunnelID int
	EncapType    string // "udp" or "ip"
	Local        string
	Remote       string
	LocalPort    int
	RemotePort   int
}

type L2tpSession struct {
	TunnelID      int
	SessionID     int
	PeerSessionID int
	Name          string
	Cookie        string
	PeerCookie    string
	L2specType    string
}

func (c *Client) L2tpAddTunnel(t *L2tpTunnel) error {
	args := []string{"l2tp", "add", "tunnel",
		"tunnel_id", fmt.Sprintf("%d", t.TunnelID),
		"peer_tunnel_id", fmt.Sprintf("%d", t.PeerTunnelID),
		"encap", t.EncapType,
		"local", t.Local,
		"remote", t.Remote,
	}
	if t.EncapType == "udp" {
		args = append(args, "udp_sport", fmt.Sprintf("%d", t.LocalPort))
		args = append(args, "udp_dport", fmt.Sprintf("%d", t.RemotePort))
	}
	return c.ipExec(args...)
}

func (c *Client) L2tpDelTunnel(tunnelID int) error {
	return c.ipExec("l2tp", "del", "tunnel", "tunnel_id", fmt.Sprintf("%d", tunnelID))
}

func (c *Client) L2tpAddSession(s *L2tpSession) error {
	args := []string{"l2tp", "add", "session",
		"tunnel_id", fmt.Sprintf("%d", s.TunnelID),
		"session_id", fmt.Sprintf("%d", s.SessionID),
		"peer_session_id", fmt.Sprintf("%d", s.PeerSessionID),
	}
	if s.Name != "" {
		args = append(args, "name", s.Name)
	}
	if s.Cookie != "" {
		args = append(args, "cookie", s.Cookie)
	}
	if s.PeerCookie != "" {
		args = append(args, "peer_cookie", s.PeerCookie)
	}
	if s.L2specType != "" {
		args = append(args, "l2spec_type", s.L2specType)
	}
	return c.ipExec(args...)
}

func (c *Client) L2tpDelSession(tunnelID, sessionID int) error {
	return c.ipExec("l2tp", "del", "session",
		"tunnel_id", fmt.Sprintf("%d", tunnelID),
		"session_id", fmt.Sprintf("%d", sessionID))
}

func (c *Client) L2tpListTunnels() (string, error) {
	return c.ipExecOutput("l2tp", "show", "tunnel")
}

func (c *Client) L2tpListSessions() (string, error) {
	return c.ipExecOutput("l2tp", "show", "session")
}
