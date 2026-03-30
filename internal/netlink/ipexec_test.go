package netlink_test

import (
	"fmt"
	"testing"
	"time"

	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
)

func TestRunInNamespace(t *testing.T) {
	name := fmt.Sprintf("tf-rns-%d", time.Now().UnixNano()%1000000)
	err := netlinkPkg.CreateNamespace(name)
	if err != nil {
		t.Fatalf("CreateNamespace failed: %v", err)
	}
	defer netlinkPkg.DeleteNamespace(name) //nolint:errcheck

	client, err := netlinkPkg.NewClient(name)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// RunInNamespace should execute function
	executed := false
	err = client.RunInNamespace(func() error {
		executed = true
		return nil
	})
	if err != nil {
		t.Fatalf("RunInNamespace failed: %v", err)
	}
	if !executed {
		t.Error("function was not executed")
	}
}

func TestRunInNamespace_error(t *testing.T) {
	name := fmt.Sprintf("tf-rnse-%d", time.Now().UnixNano()%1000000)
	err := netlinkPkg.CreateNamespace(name)
	if err != nil {
		t.Fatalf("CreateNamespace failed: %v", err)
	}
	defer netlinkPkg.DeleteNamespace(name) //nolint:errcheck

	client, err := netlinkPkg.NewClient(name)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Error should propagate
	expectedErr := fmt.Errorf("test error")
	err = client.RunInNamespace(func() error {
		return expectedErr
	})
	if err != expectedErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestRunInNamespace_defaultNamespace(t *testing.T) {
	client, err := netlinkPkg.NewClient("")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Should execute directly in default namespace
	executed := false
	err = client.RunInNamespace(func() error {
		executed = true
		return nil
	})
	if err != nil {
		t.Fatalf("RunInNamespace failed: %v", err)
	}
	if !executed {
		t.Error("function was not executed")
	}
}
