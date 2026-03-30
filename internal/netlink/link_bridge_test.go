package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
)

func TestLinkAdd_bridge(t *testing.T) {
	client := createTestClient(t)

	br := &vnl.Bridge{
		LinkAttrs: vnl.LinkAttrs{Name: "test-br0"},
	}
	if err := client.LinkAdd(br); err != nil {
		t.Fatalf("LinkAdd(bridge) failed: %v", err)
	}

	link, err := client.LinkByName("test-br0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	if link.Type() != "bridge" {
		t.Errorf("expected bridge type, got %s", link.Type())
	}
}

func TestLinkSetMaster(t *testing.T) {
	client := createTestClient(t)

	br := &vnl.Bridge{
		LinkAttrs: vnl.LinkAttrs{Name: "test-brm0"},
	}
	if err := client.LinkAdd(br); err != nil {
		t.Fatalf("LinkAdd(bridge) failed: %v", err)
	}

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-brmport"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd(dummy) failed: %v", err)
	}

	brLink, _ := client.LinkByName("test-brm0")
	port, _ := client.LinkByName("test-brmport")

	if err := client.LinkSetMaster(port, brLink); err != nil {
		t.Fatalf("LinkSetMaster failed: %v", err)
	}

	port, _ = client.LinkByName("test-brmport")
	if port.Attrs().MasterIndex != brLink.Attrs().Index {
		t.Errorf("expected master index %d, got %d", brLink.Attrs().Index, port.Attrs().MasterIndex)
	}
}

func TestLinkSetNoMaster(t *testing.T) {
	client := createTestClient(t)

	br := &vnl.Bridge{
		LinkAttrs: vnl.LinkAttrs{Name: "test-brnm0"},
	}
	if err := client.LinkAdd(br); err != nil {
		t.Fatalf("LinkAdd(bridge) failed: %v", err)
	}
	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-brnmp"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd(dummy) failed: %v", err)
	}

	brLink, _ := client.LinkByName("test-brnm0")
	port, _ := client.LinkByName("test-brnmp")
	_ = client.LinkSetMaster(port, brLink)

	if err := client.LinkSetNoMaster(port); err != nil {
		t.Fatalf("LinkSetNoMaster failed: %v", err)
	}
}

func TestLinkSetHardwareAddr(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-mac0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, _ := client.LinkByName("test-mac0")
	mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	if err := client.LinkSetHardwareAddr(link, mac); err != nil {
		t.Fatalf("LinkSetHardwareAddr failed: %v", err)
	}

	link, _ = client.LinkByName("test-mac0")
	if link.Attrs().HardwareAddr.String() != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("expected aa:bb:cc:dd:ee:ff, got %s", link.Attrs().HardwareAddr.String())
	}
}

func TestLinkSetAlias(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-alias0"},
	}
	if err := client.LinkAdd(dummy); err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, _ := client.LinkByName("test-alias0")
	if err := client.LinkSetAlias(link, "my-test-alias"); err != nil {
		t.Fatalf("LinkSetAlias failed: %v", err)
	}
}

func TestLinkAdd_veth(t *testing.T) {
	client := createTestClient(t)

	veth := &vnl.Veth{
		LinkAttrs: vnl.LinkAttrs{Name: "test-veth0"},
		PeerName:  "test-veth1",
	}
	if err := client.LinkAdd(veth); err != nil {
		t.Fatalf("LinkAdd(veth) failed: %v", err)
	}

	// Both ends should exist
	_, err := client.LinkByName("test-veth0")
	if err != nil {
		t.Fatalf("veth0 not found: %v", err)
	}
	_, err = client.LinkByName("test-veth1")
	if err != nil {
		t.Fatalf("veth1 not found: %v", err)
	}
}
