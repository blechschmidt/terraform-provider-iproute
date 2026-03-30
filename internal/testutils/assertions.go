package testutils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

// RunIPCommand runs an ip command inside a network namespace and returns the output.
func RunIPCommand(t *testing.T, namespace string, args ...string) string {
	t.Helper()

	fullArgs := append([]string{"netns", "exec", namespace, "ip"}, args...)
	cmd := exec.Command("ip", fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ip command failed: ip %s: %v: %s", strings.Join(fullArgs, " "), err, out)
	}
	return string(out)
}

// RunIPCommandNoFail runs an ip command inside a network namespace and returns the output and error.
func RunIPCommandNoFail(namespace string, args ...string) (string, error) {
	fullArgs := append([]string{"netns", "exec", namespace, "ip"}, args...)
	cmd := exec.Command("ip", fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("ip %s: %w: %s", strings.Join(fullArgs, " "), err, out)
	}
	return string(out), nil
}

// AssertLinkExists checks that a network interface exists in the namespace.
func AssertLinkExists(t *testing.T, namespace, name string) {
	t.Helper()
	output := RunIPCommand(t, namespace, "link", "show", name)
	if !strings.Contains(output, name) {
		t.Fatalf("expected link %q to exist in namespace %q, but it doesn't. Output: %s", name, namespace, output)
	}
}

// AssertLinkNotExists checks that a network interface does not exist in the namespace.
func AssertLinkNotExists(t *testing.T, namespace, name string) {
	t.Helper()
	_, err := RunIPCommandNoFail(namespace, "link", "show", name)
	if err == nil {
		t.Fatalf("expected link %q to not exist in namespace %q, but it does", name, namespace)
	}
}

// AssertAddressExists checks that an address exists on an interface in the namespace.
func AssertAddressExists(t *testing.T, namespace, iface, address string) {
	t.Helper()
	output := RunIPCommand(t, namespace, "addr", "show", "dev", iface)
	if !strings.Contains(output, address) {
		t.Fatalf("expected address %q on %q in namespace %q. Output: %s", address, iface, namespace, output)
	}
}

// AssertRouteExists checks that a route exists in the namespace.
func AssertRouteExists(t *testing.T, namespace, destination string) {
	t.Helper()
	output := RunIPCommand(t, namespace, "route", "show")
	if !strings.Contains(output, destination) {
		t.Fatalf("expected route to %q in namespace %q. Output: %s", destination, namespace, output)
	}
}

// AssertNeighborExists checks that a neighbor entry exists in the namespace.
func AssertNeighborExists(t *testing.T, namespace, address string) {
	t.Helper()
	output := RunIPCommand(t, namespace, "neigh", "show")
	if !strings.Contains(output, address) {
		t.Fatalf("expected neighbor %q in namespace %q. Output: %s", address, namespace, output)
	}
}

// AssertRuleExists checks that a routing rule exists in the namespace.
func AssertRuleExists(t *testing.T, namespace, match string) {
	t.Helper()
	output := RunIPCommand(t, namespace, "rule", "show")
	if !strings.Contains(output, match) {
		t.Fatalf("expected rule matching %q in namespace %q. Output: %s", match, namespace, output)
	}
}

// ExpectErrorRegex returns a regex for use with resource.TestStep ExpectError.
func ExpectErrorRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
