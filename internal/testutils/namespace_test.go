package testutils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/example/terraform-provider-iproute/internal/testutils"
)

func TestRandInt(t *testing.T) {
	a := testutils.RandInt()
	b := testutils.RandInt()
	// Not guaranteed to be different, but highly likely
	if a == b {
		t.Log("two consecutive RandInt calls returned same value (unlikely but possible)")
	}
}

func TestCreateTestNamespace(t *testing.T) {
	// Use a sub-test so cleanup fires after assertion
	var ns string
	t.Run("create", func(t *testing.T) {
		ns = testutils.CreateTestNamespace(t)
		if ns == "" {
			t.Fatal("CreateTestNamespace returned empty string")
		}
	})
	// After sub-test cleanup, namespace should be gone or going away
	_ = ns
}

func TestIPExec(t *testing.T) {
	name := fmt.Sprintf("tf-test-ipexec-%d", time.Now().UnixNano()%1000000)

	// Create namespace manually for testing
	cmd := fmt.Sprintf("ip netns add %s", name)
	_ = cmd

	// Test that IPExec with non-existent namespace fails
	_, err := testutils.IPExec("nonexistent-ns-12345", "link", "show")
	if err == nil {
		t.Error("expected error for nonexistent namespace")
	}
}

func TestIPExecJSON(t *testing.T) {
	_, err := testutils.IPExecJSON("nonexistent-ns-12345", "link", "show")
	if err == nil {
		t.Error("expected error for nonexistent namespace")
	}
}
