package testutils

import (
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

// CreateTestNamespace creates a unique network namespace for testing and
// registers a cleanup function to delete it when the test completes.
func CreateTestNamespace(t *testing.T) string {
	t.Helper()

	name := fmt.Sprintf("tf-test-%d-%d", time.Now().UnixNano(), rand.Intn(10000))

	cmd := exec.Command("ip", "netns", "add", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to create test namespace %q: %v: %s", name, err, out)
	}

	t.Cleanup(func() {
		cmd := exec.Command("ip", "netns", "delete", name)
		_ = cmd.Run()
	})

	// Bring up loopback in the namespace
	cmd = exec.Command("ip", "netns", "exec", name, "ip", "link", "set", "lo", "up")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to bring up loopback in namespace %q: %v: %s", name, err, out)
	}

	return name
}

// RandInt returns a random integer for test naming.
func RandInt() int {
	return rand.Intn(100000)
}
