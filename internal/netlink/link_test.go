package netlink_test

import (
	"fmt"
	"testing"
	"time"

	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
	vnl "github.com/vishvananda/netlink"
)

func createTestClient(t *testing.T) *netlinkPkg.Client {
	t.Helper()
	name := fmt.Sprintf("tf-lnk-%d", time.Now().UnixNano()%1000000)
	err := netlinkPkg.CreateNamespace(name)
	if err != nil {
		t.Fatalf("CreateNamespace failed: %v", err)
	}
	t.Cleanup(func() {
		netlinkPkg.DeleteNamespace(name) //nolint:errcheck
	})

	client, err := netlinkPkg.NewClient(name)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return client
}

func TestLinkAdd_dummy(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-d0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-d0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	if link.Attrs().Name != "test-d0" {
		t.Errorf("expected test-d0, got %s", link.Attrs().Name)
	}
}

func TestLinkDel(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-del0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-del0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	err = client.LinkDel(link)
	if err != nil {
		t.Fatalf("LinkDel failed: %v", err)
	}

	_, err = client.LinkByName("test-del0")
	if err == nil {
		t.Error("expected link to be deleted")
	}
}

func TestLinkSetUp_SetDown(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-updn0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-updn0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	err = client.LinkSetUp(link)
	if err != nil {
		t.Fatalf("LinkSetUp failed: %v", err)
	}

	err = client.LinkSetDown(link)
	if err != nil {
		t.Fatalf("LinkSetDown failed: %v", err)
	}
}

func TestLinkSetName(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-rn0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-rn0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	err = client.LinkSetName(link, "test-rn1")
	if err != nil {
		t.Fatalf("LinkSetName failed: %v", err)
	}

	link2, err := client.LinkByName("test-rn1")
	if err != nil {
		t.Fatalf("LinkByName(test-rn1) failed: %v", err)
	}

	if link2.Attrs().Name != "test-rn1" {
		t.Errorf("expected test-rn1, got %s", link2.Attrs().Name)
	}
}

func TestLinkSetMTU(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-mtu0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-mtu0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	err = client.LinkSetMTU(link, 1400)
	if err != nil {
		t.Fatalf("LinkSetMTU failed: %v", err)
	}

	link, err = client.LinkByName("test-mtu0")
	if err != nil {
		t.Fatalf("LinkByName failed after MTU change: %v", err)
	}

	if link.Attrs().MTU != 1400 {
		t.Errorf("expected MTU 1400, got %d", link.Attrs().MTU)
	}
}

func TestLinkList(t *testing.T) {
	client := createTestClient(t)

	links, err := client.LinkList()
	if err != nil {
		t.Fatalf("LinkList failed: %v", err)
	}

	// Should have at least loopback
	found := false
	for _, l := range links {
		if l.Attrs().Name == "lo" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find lo in link list")
	}
}

func TestLinkSetTxQLen(t *testing.T) {
	client := createTestClient(t)

	dummy := &vnl.Dummy{
		LinkAttrs: vnl.LinkAttrs{Name: "test-txq0"},
	}
	err := client.LinkAdd(dummy)
	if err != nil {
		t.Fatalf("LinkAdd failed: %v", err)
	}

	link, err := client.LinkByName("test-txq0")
	if err != nil {
		t.Fatalf("LinkByName failed: %v", err)
	}

	err = client.LinkSetTxQLen(link, 500)
	if err != nil {
		t.Fatalf("LinkSetTxQLen failed: %v", err)
	}
}
