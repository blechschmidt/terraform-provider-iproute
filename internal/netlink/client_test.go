package netlink_test

import (
	"testing"

	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
)

func TestNewClient_defaultNamespace(t *testing.T) {
	client, err := netlinkPkg.NewClient("")
	if err != nil {
		t.Fatalf("NewClient with empty namespace should succeed: %v", err)
	}
	defer client.Close()

	if client.Namespace != "" {
		t.Errorf("expected empty namespace, got %q", client.Namespace)
	}

	// Should be able to list links
	links, err := client.LinkList()
	if err != nil {
		t.Fatalf("LinkList failed: %v", err)
	}

	// Should have at least loopback
	if len(links) == 0 {
		t.Error("expected at least one link (lo)")
	}

	found := false
	for _, l := range links {
		if l.Attrs().Name == "lo" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find loopback interface")
	}
}

func TestNewClient_invalidNamespace(t *testing.T) {
	_, err := netlinkPkg.NewClient("nonexistent-ns-test-12345")
	if err == nil {
		t.Error("expected error for nonexistent namespace")
	}
}

func TestLinkByName_loopback(t *testing.T) {
	client, err := netlinkPkg.NewClient("")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	link, err := client.LinkByName("lo")
	if err != nil {
		t.Fatalf("LinkByName(lo) failed: %v", err)
	}

	if link.Attrs().Name != "lo" {
		t.Errorf("expected lo, got %s", link.Attrs().Name)
	}
}

func TestLinkByName_notFound(t *testing.T) {
	client, err := netlinkPkg.NewClient("")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	_, err = client.LinkByName("nonexistent-iface-12345")
	if err == nil {
		t.Error("expected error for nonexistent interface")
	}
}

func TestLinkByIndex_loopback(t *testing.T) {
	client, err := netlinkPkg.NewClient("")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	link, err := client.LinkByIndex(1) // lo is always index 1
	if err != nil {
		t.Fatalf("LinkByIndex(1) failed: %v", err)
	}

	if link.Attrs().Name != "lo" {
		t.Errorf("expected lo at index 1, got %s", link.Attrs().Name)
	}
}
