package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestNeighAdd_NeighDel(t *testing.T) {
	client := createTestClient(t)

	// Create dummy + address
	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-ne0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-ne0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}
	if err := client.LinkSetUp(link); err != nil {
		t.Fatalf("LinkSetUp failed: %v", err)
	}
	addr := &vnl.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP("10.0.0.1"),
			Mask: net.CIDRMask(24, 32),
		},
	}
	if err := client.AddrAdd(link, addr); err != nil {
		t.Fatalf("AddrAdd failed: %v", err)
	}

	mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	neigh := &vnl.Neigh{
		LinkIndex:    link.Attrs().Index,
		IP:           net.ParseIP("10.0.0.2"),
		HardwareAddr: mac,
		State:        vnl.NUD_PERMANENT,
	}

	if err := client.NeighAdd(neigh); err != nil {
		t.Fatalf("NeighAdd failed: %v", err)
	}

	neighs, err := client.NeighList(link.Attrs().Index, unix.AF_INET)
	if err != nil {
		t.Fatalf("NeighList failed: %v", err)
	}

	found := false
	for _, n := range neighs {
		if n.IP.Equal(net.ParseIP("10.0.0.2")) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find neighbor 10.0.0.2")
	}

	if err := client.NeighDel(neigh); err != nil {
		t.Fatalf("NeighDel failed: %v", err)
	}
}

func TestNeighSet(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-neset0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-neset0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}
	if err := client.LinkSetUp(link); err != nil {
		t.Fatalf("LinkSetUp failed: %v", err)
	}
	addr := &vnl.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP("10.0.0.1"),
			Mask: net.CIDRMask(24, 32),
		},
	}
	if err := client.AddrAdd(link, addr); err != nil {
		t.Fatalf("AddrAdd failed: %v", err)
	}

	mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:01")
	neigh := &vnl.Neigh{
		LinkIndex:    link.Attrs().Index,
		IP:           net.ParseIP("10.0.0.3"),
		HardwareAddr: mac,
		State:        vnl.NUD_PERMANENT,
	}

	// NeighSet should work even if neighbor doesn't exist yet (it does add+set)
	if err := client.NeighSet(neigh); err != nil {
		t.Fatalf("NeighSet failed: %v", err)
	}
}

func TestNeighList_empty(t *testing.T) {
	client := createTestClient(t)

	// List on loopback (no neighbors expected)
	lo, err := client.LinkByName("lo")
	if err != nil {
		t.Fatalf("LinkByName(lo) failed: %v", err)
	}

	neighs, err := client.NeighList(lo.Attrs().Index, unix.AF_INET)
	if err != nil {
		t.Fatalf("NeighList failed: %v", err)
	}

	if len(neighs) != 0 {
		t.Errorf("expected no neighbors on lo, got %d", len(neighs))
	}
}
