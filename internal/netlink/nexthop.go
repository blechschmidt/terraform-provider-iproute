package netlink

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/vishvananda/netns"
	"golang.org/x/sys/unix"
)

// Nexthop represents a kernel nexthop object (RTM_NEWNEXTHOP).
type Nexthop struct {
	ID        int
	Gateway   net.IP
	LinkIndex int
	Blackhole bool
	Family    int
	Group     []NexthopGroupMember
	FDB       bool
}

// NexthopGroupMember represents a member of a nexthop group.
type NexthopGroupMember struct {
	ID     int
	Weight int
}

// Raw rtnetlink constants for nexthop objects.
const (
	RTM_NEWNEXTHOP = 104
	RTM_DELNEXTHOP = 105
	RTM_GETNEXTHOP = 106

	NHA_ID          = 1
	NHA_GROUP       = 2
	NHA_GROUP_TYPE  = 3
	NHA_BLACKHOLE   = 4
	NHA_OIF         = 5
	NHA_GATEWAY     = 6
	NHA_FDB         = 11

	NEXTHOP_GRP_TYPE_MPATH = 0
	NEXTHOP_GRP_TYPE_RES   = 1
)

func (c *Client) NexthopAdd(nh *Nexthop) error {
	return c.RunInNamespace(func() error {
		return nexthopModify(RTM_NEWNEXTHOP, unix.NLM_F_CREATE|unix.NLM_F_EXCL, nh)
	})
}

func (c *Client) NexthopDel(nh *Nexthop) error {
	return c.RunInNamespace(func() error {
		return nexthopModify(RTM_DELNEXTHOP, 0, nh)
	})
}

func (c *Client) NexthopList() ([]Nexthop, error) {
	var result []Nexthop
	err := c.RunInNamespace(func() error {
		var err error
		result, err = nexthopList()
		return err
	})
	return result, err
}

func (c *Client) NexthopGet(id int) (*Nexthop, error) {
	nhs, err := c.NexthopList()
	if err != nil {
		return nil, err
	}
	for _, nh := range nhs {
		if nh.ID == id {
			return &nh, nil
		}
	}
	return nil, fmt.Errorf("nexthop %d not found", id)
}

func nexthopModify(cmd, flags int, nh *Nexthop) error {
	s, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.NETLINK_ROUTE)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(s)

	lsa := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	if err := unix.Bind(s, lsa); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	family := byte(unix.AF_INET)
	if nh.Family != 0 {
		family = byte(nh.Family)
	}
	if nh.Gateway != nil && nh.Gateway.To4() == nil {
		family = byte(unix.AF_INET6)
	}

	// Build nhmsg header: family(1) + scope(1) + protocol(1) + resvd(1) + flags(4)
	nhmsg := make([]byte, 8)
	nhmsg[0] = family

	// Build attributes
	var attrs []byte

	// NHA_ID
	attrs = append(attrs, nlattr(NHA_ID, uint32Bytes(uint32(nh.ID)))...)

	if len(nh.Group) > 0 {
		// NHA_GROUP
		groupData := make([]byte, len(nh.Group)*8)
		for i, m := range nh.Group {
			binary.LittleEndian.PutUint32(groupData[i*8:], uint32(m.ID))
			groupData[i*8+4] = byte(m.Weight)
		}
		attrs = append(attrs, nlattr(NHA_GROUP, groupData)...)
	} else if nh.Blackhole {
		attrs = append(attrs, nlattr(NHA_BLACKHOLE, nil)...)
	} else {
		if nh.LinkIndex > 0 {
			attrs = append(attrs, nlattr(NHA_OIF, uint32Bytes(uint32(nh.LinkIndex)))...)
		}
		if nh.Gateway != nil {
			gw := nh.Gateway.To4()
			if gw == nil {
				gw = nh.Gateway.To16()
			}
			attrs = append(attrs, nlattr(NHA_GATEWAY, gw)...)
		}
	}

	if nh.FDB {
		attrs = append(attrs, nlattr(NHA_FDB, nil)...)
	}

	payload := append(nhmsg, attrs...)

	msg := buildNetlinkMessage(uint16(cmd), uint16(unix.NLM_F_REQUEST|unix.NLM_F_ACK|flags), payload)

	if err := unix.Sendto(s, msg, 0, lsa); err != nil {
		return fmt.Errorf("sendto: %w", err)
	}

	return recvAck(s)
}

func nexthopList() ([]Nexthop, error) {
	s, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.NETLINK_ROUTE)
	if err != nil {
		return nil, fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(s)

	lsa := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	if err := unix.Bind(s, lsa); err != nil {
		return nil, fmt.Errorf("bind: %w", err)
	}

	nhmsg := make([]byte, 8) // family=0 (AF_UNSPEC)
	msg := buildNetlinkMessage(RTM_GETNEXTHOP, unix.NLM_F_REQUEST|unix.NLM_F_DUMP, nhmsg)

	if err := unix.Sendto(s, msg, 0, lsa); err != nil {
		return nil, fmt.Errorf("sendto: %w", err)
	}

	var result []Nexthop
	buf := make([]byte, 65536)
	for {
		n, _, err := unix.Recvfrom(s, buf, 0)
		if err != nil {
			return nil, fmt.Errorf("recvfrom: %w", err)
		}
		data := buf[:n]
		for len(data) >= 16 {
			msgLen := binary.LittleEndian.Uint32(data[0:4])
			msgType := binary.LittleEndian.Uint16(data[4:6])
			if msgType == unix.NLMSG_DONE {
				return result, nil
			}
			if msgType == unix.NLMSG_ERROR {
				errno := int32(binary.LittleEndian.Uint32(data[16:20]))
				if errno != 0 {
					return nil, fmt.Errorf("netlink error: %d", errno)
				}
				return result, nil
			}
			if msgType == RTM_NEWNEXTHOP && msgLen > 24 {
				nh := parseNexthop(data[16:msgLen], data[16]) // pass family byte
				result = append(result, nh)
			}
			data = data[align4(int(msgLen)):]
		}
	}
}

func parseNexthop(data []byte, family byte) Nexthop {
	nh := Nexthop{Family: int(family)}
	if len(data) < 8 {
		return nh
	}
	// Skip nhmsg header (8 bytes)
	attrs := data[8:]
	for len(attrs) >= 4 {
		attrLen := int(binary.LittleEndian.Uint16(attrs[0:2]))
		attrType := int(binary.LittleEndian.Uint16(attrs[2:4]))
		if attrLen < 4 || attrLen > len(attrs) {
			break
		}
		payload := attrs[4:attrLen]
		switch attrType {
		case NHA_ID:
			if len(payload) >= 4 {
				nh.ID = int(binary.LittleEndian.Uint32(payload))
			}
		case NHA_OIF:
			if len(payload) >= 4 {
				nh.LinkIndex = int(binary.LittleEndian.Uint32(payload))
			}
		case NHA_GATEWAY:
			nh.Gateway = make(net.IP, len(payload))
			copy(nh.Gateway, payload)
		case NHA_BLACKHOLE:
			nh.Blackhole = true
		case NHA_GROUP:
			for i := 0; i+8 <= len(payload); i += 8 {
				m := NexthopGroupMember{
					ID:     int(binary.LittleEndian.Uint32(payload[i:])),
					Weight: int(payload[i+4]),
				}
				nh.Group = append(nh.Group, m)
			}
		case NHA_FDB:
			nh.FDB = true
		}
		attrs = attrs[align4(attrLen):]
	}
	return nh
}

// helper functions for building raw netlink messages

func nlattr(typ int, data []byte) []byte {
	l := 4 + len(data)
	b := make([]byte, align4(l))
	binary.LittleEndian.PutUint16(b[0:2], uint16(l))
	binary.LittleEndian.PutUint16(b[2:4], uint16(typ))
	copy(b[4:], data)
	return b
}

func uint32Bytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

func buildNetlinkMessage(msgType, flags uint16, payload []byte) []byte {
	msgLen := 16 + len(payload)
	msg := make([]byte, align4(msgLen))
	binary.LittleEndian.PutUint32(msg[0:4], uint32(msgLen))
	binary.LittleEndian.PutUint16(msg[4:6], msgType)
	binary.LittleEndian.PutUint16(msg[6:8], flags)
	binary.LittleEndian.PutUint32(msg[8:12], 1) // seq
	binary.LittleEndian.PutUint32(msg[12:16], 0) // pid
	copy(msg[16:], payload)
	return msg
}

func recvAck(s int) error {
	buf := make([]byte, 4096)
	n, _, err := unix.Recvfrom(s, buf, 0)
	if err != nil {
		return fmt.Errorf("recvfrom: %w", err)
	}
	if n < 20 {
		return fmt.Errorf("short ack message")
	}
	msgType := binary.LittleEndian.Uint16(buf[4:6])
	if msgType == unix.NLMSG_ERROR {
		errno := int32(binary.LittleEndian.Uint32(buf[16:20]))
		if errno != 0 {
			return fmt.Errorf("netlink error: %s", unix.Errno(-errno))
		}
	}
	return nil
}

func align4(v int) int {
	return (v + 3) &^ 3
}

// Ensure netns import is used (for RunInNamespace via client.go)
var _ = netns.None
