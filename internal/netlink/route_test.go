package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestRouteAdd_RouteDel(t *testing.T) {
	client := createTestClient(t)

	// Create dummy interface and address
	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-rt0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-rt0")
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

	_, dst, _ := net.ParseCIDR("10.1.0.0/24")
	route := &vnl.Route{
		Dst:       dst,
		Gw:        net.ParseIP("10.0.0.254"),
		LinkIndex: link.Attrs().Index,
	}

	if err := client.RouteAdd(route); err != nil {
		t.Fatalf("RouteAdd failed: %v", err)
	}

	routes, err := client.RouteList(nil, unix.AF_INET)
	if err != nil {
		t.Fatalf("RouteList failed: %v", err)
	}

	found := false
	for _, r := range routes {
		if r.Dst != nil && r.Dst.String() == "10.1.0.0/24" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find route 10.1.0.0/24")
	}

	if err := client.RouteDel(route); err != nil {
		t.Fatalf("RouteDel failed: %v", err)
	}
}

func TestRouteReplace(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-rtrpl0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-rtrpl0")
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

	_, dst, _ := net.ParseCIDR("10.2.0.0/24")
	route := &vnl.Route{
		Dst:       dst,
		Gw:        net.ParseIP("10.0.0.254"),
		LinkIndex: link.Attrs().Index,
	}

	if err := client.RouteAdd(route); err != nil {
		t.Fatalf("RouteAdd failed: %v", err)
	}

	// Replace with same route should succeed
	if err := client.RouteReplace(route); err != nil {
		t.Fatalf("RouteReplace failed: %v", err)
	}
}
