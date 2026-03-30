package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestAddrAdd_AddrList(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-a0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-a0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
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

	addrs, err := client.AddrList(link, unix.AF_INET)
	if err != nil {
		t.Fatalf("AddrList failed: %v", err)
	}

	found := false
	for _, a := range addrs {
		if a.IPNet.IP.Equal(net.ParseIP("10.0.0.1")) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 10.0.0.1 in address list")
	}
}

func TestAddrDel(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-adel0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-adel0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	addr := &vnl.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP("10.0.1.1"),
			Mask: net.CIDRMask(24, 32),
		},
	}
	if err := client.AddrAdd(link, addr); err != nil {
		t.Fatalf("AddrAdd failed: %v", err)
	}

	if err := client.AddrDel(link, addr); err != nil {
		t.Fatalf("AddrDel failed: %v", err)
	}

	addrs, err := client.AddrList(link, unix.AF_INET)
	if err != nil {
		t.Fatalf("AddrList failed: %v", err)
	}

	for _, a := range addrs {
		if a.IPNet.IP.Equal(net.ParseIP("10.0.1.1")) {
			t.Error("expected address 10.0.1.1 to be deleted")
		}
	}
}

func TestAddrReplace(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-arpl0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}
	link, err := client.LinkByName("test-arpl0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	addr := &vnl.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP("10.0.2.1"),
			Mask: net.CIDRMask(24, 32),
		},
	}
	if err := client.AddrAdd(link, addr); err != nil {
		t.Fatalf("AddrAdd failed: %v", err)
	}

	// Replace should work without error
	if err := client.AddrReplace(link, addr); err != nil {
		t.Fatalf("AddrReplace failed: %v", err)
	}
}
