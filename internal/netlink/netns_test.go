package netlink_test

import (
	"fmt"
	"testing"
	"time"

	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
)

func TestCreateDeleteNamespace(t *testing.T) {
	name := fmt.Sprintf("tf-test-%d", time.Now().UnixNano())

	err := netlinkPkg.CreateNamespace(name)
	if err != nil {
		t.Fatalf("CreateNamespace(%q) failed: %v", name, err)
	}
	defer netlinkPkg.DeleteNamespace(name) //nolint:errcheck

	exists, err := netlinkPkg.NamespaceExists(name)
	if err != nil {
		t.Fatalf("NamespaceExists(%q) failed: %v", name, err)
	}
	if !exists {
		t.Errorf("expected namespace %q to exist", name)
	}

	err = netlinkPkg.DeleteNamespace(name)
	if err != nil {
		t.Fatalf("DeleteNamespace(%q) failed: %v", name, err)
	}

	exists, err = netlinkPkg.NamespaceExists(name)
	if err != nil {
		t.Fatalf("NamespaceExists(%q) failed: %v", name, err)
	}
	if exists {
		t.Errorf("expected namespace %q to not exist", name)
	}
}

func TestNamespaceExists_nonexistent(t *testing.T) {
	exists, err := netlinkPkg.NamespaceExists("nonexistent-ns-99999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected namespace to not exist")
	}
}

func TestListNamespaces(t *testing.T) {
	// List should work even if empty
	_, err := netlinkPkg.ListNamespaces()
	if err != nil {
		t.Fatalf("ListNamespaces failed: %v", err)
	}
}

func TestCreateNamespace_duplicate(t *testing.T) {
	name := fmt.Sprintf("tf-test-dup-%d", time.Now().UnixNano())

	err := netlinkPkg.CreateNamespace(name)
	if err != nil {
		t.Fatalf("CreateNamespace(%q) failed: %v", name, err)
	}
	defer netlinkPkg.DeleteNamespace(name) //nolint:errcheck

	// Creating again should fail or be idempotent
	err = netlinkPkg.CreateNamespace(name)
	if err == nil {
		// Some systems allow re-creating, which is fine
		t.Log("duplicate namespace creation succeeded (idempotent)")
	}
}
