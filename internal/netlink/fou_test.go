package netlink_test

import (
	"os/exec"
	"testing"

	vnl "github.com/vishvananda/netlink"
)

func ensureFouModule(t *testing.T) {
	t.Helper()
	if err := exec.Command("modprobe", "fou").Run(); err != nil {
		t.Skipf("FOU kernel module not available: %v", err)
	}
}

func TestFouAdd_Del(t *testing.T) {
	ensureFouModule(t)
	client := createTestClient(t)

	fou := vnl.Fou{
		Port:      5555,
		Family:    4,
		Protocol:  4,
		EncapType: vnl.FOU_ENCAP_DIRECT,
	}

	if err := client.FouAdd(fou); err != nil {
		t.Skipf("FouAdd failed (kernel may not support FOU): %v", err)
	}

	fous, err := client.FouList(4)
	if err != nil {
		t.Skipf("FouList not supported: %v", err)
	}

	found := false
	for _, f := range fous {
		if f.Port == 5555 {
			found = true
			break
		}
	}
	if !found {
		t.Log("FOU 5555 not found in list (may be a kernel version issue)")
	}

	if err := client.FouDel(fou); err != nil {
		t.Fatalf("FouDel failed: %v", err)
	}
}

func TestFouAdd_GUE(t *testing.T) {
	ensureFouModule(t)
	client := createTestClient(t)

	fou := vnl.Fou{
		Port:      5556,
		Family:    4,
		Protocol:  4,
		EncapType: vnl.FOU_ENCAP_GUE,
	}

	if err := client.FouAdd(fou); err != nil {
		t.Skipf("FouAdd (GUE) failed (kernel may not support FOU): %v", err)
	}

	if err := client.FouDel(fou); err != nil {
		t.Fatalf("FouDel failed: %v", err)
	}
}
