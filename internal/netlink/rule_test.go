package netlink_test

import (
	"net"
	"testing"

	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestRuleAdd_RuleDel(t *testing.T) {
	client := createTestClient(t)

	_, src, _ := net.ParseCIDR("10.99.0.0/24")
	rule := vnl.NewRule()
	rule.Src = src
	rule.Table = 999
	rule.Priority = 9999

	if err := client.RuleAdd(rule); err != nil {
		t.Fatalf("RuleAdd failed: %v", err)
	}

	rules, err := client.RuleList(unix.AF_INET)
	if err != nil {
		t.Fatalf("RuleList failed: %v", err)
	}

	found := false
	for _, r := range rules {
		if r.Table == 999 && r.Priority == 9999 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find rule with table 999")
	}

	if err := client.RuleDel(rule); err != nil {
		t.Fatalf("RuleDel failed: %v", err)
	}
}

func TestRuleList(t *testing.T) {
	client := createTestClient(t)

	// Every namespace should have default rules
	rules, err := client.RuleList(unix.AF_INET)
	if err != nil {
		t.Fatalf("RuleList failed: %v", err)
	}

	if len(rules) == 0 {
		t.Error("expected at least one rule (default)")
	}
}
