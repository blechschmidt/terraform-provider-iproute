package netlink

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const netnsPath = "/run/netns"

func CreateNamespace(name string) error {
	cmd := exec.Command("ip", "netns", "add", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create namespace %q: %w: %s", name, err, out)
	}
	return nil
}

func DeleteNamespace(name string) error {
	cmd := exec.Command("ip", "netns", "delete", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete namespace %q: %w: %s", name, err, out)
	}
	return nil
}

func NamespaceExists(name string) (bool, error) {
	_, err := os.Stat(filepath.Join(netnsPath, name))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ListNamespaces() ([]string, error) {
	entries, err := os.ReadDir(netnsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}
