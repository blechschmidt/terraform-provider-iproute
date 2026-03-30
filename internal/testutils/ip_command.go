package testutils

import (
	"os/exec"
	"strings"
)

// IPExec runs an ip command in a namespace and returns stdout/stderr and any error.
func IPExec(namespace string, args ...string) (string, error) {
	fullArgs := append([]string{"netns", "exec", namespace, "ip"}, args...)
	cmd := exec.Command("ip", fullArgs...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// IPExecJSON runs an ip command with -j (JSON) in a namespace.
func IPExecJSON(namespace string, args ...string) (string, error) {
	jsonArgs := append([]string{"-j"}, args...)
	return IPExec(namespace, jsonArgs...)
}
