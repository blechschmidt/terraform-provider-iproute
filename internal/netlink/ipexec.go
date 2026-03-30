package netlink

import (
	"fmt"
	"os/exec"
	"strings"
)

// ipExec runs an ip command in the client's namespace.
func (c *Client) ipExec(args ...string) error {
	var cmd *exec.Cmd
	if c.Namespace != "" {
		fullArgs := append([]string{"netns", "exec", c.Namespace, "ip"}, args...)
		cmd = exec.Command("ip", fullArgs...)
	} else {
		cmd = exec.Command("ip", args...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ip %s: %w: %s", strings.Join(args, " "), err, out)
	}
	return nil
}

// ipExecOutput runs an ip command in the client's namespace and returns output.
func (c *Client) ipExecOutput(args ...string) (string, error) {
	var cmd *exec.Cmd
	if c.Namespace != "" {
		fullArgs := append([]string{"netns", "exec", c.Namespace, "ip"}, args...)
		cmd = exec.Command("ip", fullArgs...)
	} else {
		cmd = exec.Command("ip", args...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ip %s: %w: %s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
