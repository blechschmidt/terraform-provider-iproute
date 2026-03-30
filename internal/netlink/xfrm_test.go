package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
)

func TestXfrmStateAdd_Del(t *testing.T) {
	client := createTestClient(t)

	state := &vnl.XfrmState{
		Src:   net.ParseIP("10.0.0.1"),
		Dst:   net.ParseIP("10.0.0.2"),
		Proto: vnl.XFRM_PROTO_ESP,
		Spi:   1234,
		Mode:  vnl.XFRM_MODE_TUNNEL,
		Aead: &vnl.XfrmStateAlgo{
			Name:   "rfc4106(gcm(aes))",
			Key:    []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67},
			ICVLen: 128,
		},
	}

	if err := client.XfrmStateAdd(state); err != nil {
		t.Fatalf("XfrmStateAdd failed: %v", err)
	}

	states, err := client.XfrmStateList(0)
	if err != nil {
		t.Fatalf("XfrmStateList failed: %v", err)
	}

	found := false
	for _, s := range states {
		if s.Spi == 1234 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find xfrm state with SPI 1234")
	}

	if err := client.XfrmStateDel(state); err != nil {
		t.Fatalf("XfrmStateDel failed: %v", err)
	}
}

func TestXfrmPolicyAdd_Del(t *testing.T) {
	client := createTestClient(t)

	_, src, _ := net.ParseCIDR("10.0.0.0/24")
	_, dst, _ := net.ParseCIDR("10.0.1.0/24")

	policy := &vnl.XfrmPolicy{
		Src: src,
		Dst: dst,
		Dir: vnl.XFRM_DIR_OUT,
		Tmpls: []vnl.XfrmPolicyTmpl{
			{
				Src:   net.ParseIP("10.0.0.1"),
				Dst:   net.ParseIP("10.0.0.2"),
				Proto: vnl.XFRM_PROTO_ESP,
				Mode:  vnl.XFRM_MODE_TUNNEL,
			},
		},
	}

	if err := client.XfrmPolicyAdd(policy); err != nil {
		t.Fatalf("XfrmPolicyAdd failed: %v", err)
	}

	policies, err := client.XfrmPolicyList(0)
	if err != nil {
		t.Fatalf("XfrmPolicyList failed: %v", err)
	}

	found := false
	for _, p := range policies {
		if p.Src.String() == "10.0.0.0/24" && p.Dst.String() == "10.0.1.0/24" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find xfrm policy")
	}

	if err := client.XfrmPolicyDel(policy); err != nil {
		t.Fatalf("XfrmPolicyDel failed: %v", err)
	}
}

func TestXfrmStateList_empty(t *testing.T) {
	client := createTestClient(t)

	states, err := client.XfrmStateList(0)
	if err != nil {
		t.Fatalf("XfrmStateList failed: %v", err)
	}

	if len(states) != 0 {
		t.Errorf("expected no xfrm states, got %d", len(states))
	}
}

func TestXfrmPolicyList_empty(t *testing.T) {
	client := createTestClient(t)

	policies, err := client.XfrmPolicyList(0)
	if err != nil {
		t.Fatalf("XfrmPolicyList failed: %v", err)
	}

	if len(policies) != 0 {
		t.Errorf("expected no xfrm policies, got %d", len(policies))
	}
}
