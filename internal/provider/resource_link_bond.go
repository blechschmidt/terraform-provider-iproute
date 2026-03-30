package provider

import (
	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func bondSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Bond-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"mode":             schema.StringAttribute{Optional: true, Description: "Bonding mode (balance-rr, active-backup, balance-xor, broadcast, 802.3ad, balance-tlb, balance-alb)."},
			"miimon":           schema.Int64Attribute{Optional: true, Description: "MII monitoring interval (ms)."},
			"up_delay":         schema.Int64Attribute{Optional: true, Description: "Delay before enabling slave (ms)."},
			"down_delay":       schema.Int64Attribute{Optional: true, Description: "Delay before disabling slave (ms)."},
			"primary":          schema.StringAttribute{Optional: true, Description: "Primary slave interface."},
			"lacp_rate":        schema.StringAttribute{Optional: true, Description: "LACP rate (slow or fast)."},
			"xmit_hash_policy": schema.StringAttribute{Optional: true, Description: "Transmit hash policy."},
		},
	}
}

var bondModeMap = map[string]vnl.BondMode{
	"balance-rr":    vnl.BOND_MODE_BALANCE_RR,
	"active-backup": vnl.BOND_MODE_ACTIVE_BACKUP,
	"balance-xor":   vnl.BOND_MODE_BALANCE_XOR,
	"broadcast":     vnl.BOND_MODE_BROADCAST,
	"802.3ad":       vnl.BOND_MODE_802_3AD,
	"balance-tlb":   vnl.BOND_MODE_BALANCE_TLB,
	"balance-alb":   vnl.BOND_MODE_BALANCE_ALB,
}

func buildBond(name string, cfg *models.BondConfig) (*vnl.Bond, error) {
	mode := vnl.BOND_MODE_BALANCE_RR
	if cfg != nil && !cfg.Mode.IsNull() {
		if m, ok := bondModeMap[cfg.Mode.ValueString()]; ok {
			mode = m
		}
	}
	bond := vnl.NewLinkBond(vnl.LinkAttrs{Name: name})
	bond.Mode = mode
	if cfg != nil {
		if !cfg.MiiMon.IsNull() {
			v := int(cfg.MiiMon.ValueInt64())
			bond.Miimon = v
		}
		if !cfg.UpDelay.IsNull() {
			v := int(cfg.UpDelay.ValueInt64())
			bond.UpDelay = v
		}
		if !cfg.DownDelay.IsNull() {
			v := int(cfg.DownDelay.ValueInt64())
			bond.DownDelay = v
		}
		if !cfg.LacpRate.IsNull() {
			if cfg.LacpRate.ValueString() == "fast" {
				bond.LacpRate = vnl.BOND_LACP_RATE_FAST
			} else {
				bond.LacpRate = vnl.BOND_LACP_RATE_SLOW
			}
		}
		if !cfg.XmitHashPolicy.IsNull() {
			switch cfg.XmitHashPolicy.ValueString() {
			case "layer2":
				bond.XmitHashPolicy = vnl.BOND_XMIT_HASH_POLICY_LAYER2
			case "layer3+4":
				bond.XmitHashPolicy = vnl.BOND_XMIT_HASH_POLICY_LAYER3_4
			case "layer2+3":
				bond.XmitHashPolicy = vnl.BOND_XMIT_HASH_POLICY_LAYER2_3
			case "encap2+3":
				bond.XmitHashPolicy = vnl.BOND_XMIT_HASH_POLICY_ENCAP2_3
			case "encap3+4":
				bond.XmitHashPolicy = vnl.BOND_XMIT_HASH_POLICY_ENCAP3_4
			}
		}
	}
	return bond, nil
}
