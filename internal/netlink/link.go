package netlink

import (
	"fmt"
	"net"

	vnl "github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

// iflaBrGroupFwdMask is IFLA_BR_GROUP_FWD_MASK from <linux/if_link.h>.
const iflaBrGroupFwdMask = 9

// LinkSetBridgeGroupFwdMask sets a bridge's group_fwd_mask via a raw RTM_NEWLINK
// netlink message (IFLA_BR_GROUP_FWD_MASK), which netlink v1.3.0 does not model
// on its Bridge type. The request is executed from a thread switched into the
// client's namespace so the netlink socket is opened there; this avoids sysfs,
// whose per-namespace /sys view is unavailable for a foreign namespace. No
// shell is involved.
func (c *Client) LinkSetBridgeGroupFwdMask(name string, mask uint16) error {
	link, err := c.Handle.LinkByName(name)
	if err != nil {
		return fmt.Errorf("looking up bridge %q: %w", name, err)
	}
	index := link.Attrs().Index

	return c.RunInNamespace(func() error {
		req := nl.NewNetlinkRequest(unix.RTM_NEWLINK, unix.NLM_F_ACK)
		msg := nl.NewIfInfomsg(unix.AF_UNSPEC)
		msg.Index = int32(index)
		req.AddData(msg)

		linkInfo := nl.NewRtAttr(unix.IFLA_LINKINFO, nil)
		linkInfo.AddRtAttr(nl.IFLA_INFO_KIND, nl.NonZeroTerminated("bridge"))
		data := linkInfo.AddRtAttr(nl.IFLA_INFO_DATA, nil)
		data.AddRtAttr(iflaBrGroupFwdMask, nl.Uint16Attr(mask))
		req.AddData(linkInfo)

		if _, err := req.Execute(unix.NETLINK_ROUTE, 0); err != nil {
			return fmt.Errorf("setting group_fwd_mask on %q: %w", name, err)
		}
		return nil
	})
}

func (c *Client) LinkAdd(link vnl.Link) error {
	// TUN/TAP devices are created via /dev/net/tun ioctl, not netlink messages.
	// The Handle's namespace scoping only affects netlink sockets, so we must
	// switch the OS thread to the target namespace for tuntap creation.
	if _, ok := link.(*vnl.Tuntap); ok {
		return c.RunInNamespace(func() error {
			return c.Handle.LinkAdd(link)
		})
	}
	return c.Handle.LinkAdd(link)
}

func (c *Client) LinkDel(link vnl.Link) error {
	return c.Handle.LinkDel(link)
}

func (c *Client) LinkSetUp(link vnl.Link) error {
	return c.Handle.LinkSetUp(link)
}

func (c *Client) LinkSetDown(link vnl.Link) error {
	return c.Handle.LinkSetDown(link)
}

func (c *Client) LinkSetName(link vnl.Link, name string) error {
	return c.Handle.LinkSetName(link, name)
}

func (c *Client) LinkSetMTU(link vnl.Link, mtu int) error {
	return c.Handle.LinkSetMTU(link, mtu)
}

func (c *Client) LinkSetMaster(link vnl.Link, master vnl.Link) error {
	return c.Handle.LinkSetMaster(link, master)
}

func (c *Client) LinkSetNoMaster(link vnl.Link) error {
	return c.Handle.LinkSetNoMaster(link)
}

func (c *Client) LinkSetHardwareAddr(link vnl.Link, addr net.HardwareAddr) error {
	return c.Handle.LinkSetHardwareAddr(link, addr)
}

func (c *Client) LinkSetTxQLen(link vnl.Link, qlen int) error {
	return c.Handle.LinkSetTxQLen(link, qlen)
}

func (c *Client) LinkSetAlias(link vnl.Link, alias string) error {
	return c.Handle.LinkSetAlias(link, alias)
}

func (c *Client) LinkByName(name string) (vnl.Link, error) {
	return c.Handle.LinkByName(name)
}

func (c *Client) LinkByIndex(index int) (vnl.Link, error) {
	return c.Handle.LinkByIndex(index)
}

func (c *Client) LinkList() ([]vnl.Link, error) {
	return c.Handle.LinkList()
}
