package provider

import (
	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func bridgeSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Bridge-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"stp":            schema.BoolAttribute{Optional: true, Description: "Enable STP."},
			"hello_time":     schema.Int64Attribute{Optional: true, Description: "STP hello time (centiseconds)."},
			"max_age":        schema.Int64Attribute{Optional: true, Description: "STP max age (centiseconds)."},
			"forward_delay":  schema.Int64Attribute{Optional: true, Description: "STP forward delay (centiseconds)."},
			"vlan_filtering": schema.BoolAttribute{Optional: true, Description: "Enable VLAN filtering."},
			"default_pvid":   schema.Int64Attribute{Optional: true, Description: "Default PVID."},
			"ageing_time":    schema.Int64Attribute{Optional: true, Description: "MAC ageing time (centiseconds)."},
			"group_fwd_mask": schema.Int64Attribute{
				Optional: true,
				Description: "Per-bridge group forwarding mask (IFLA_BR_GROUP_FWD_MASK). " +
					"A bitmap of the 802.1D reserved link-local multicast addresses " +
					"(01:80:c2:00:00:0X) the bridge forwards instead of dropping. " +
					"For example 0x4000 (16384) forwards LLDP (01:80:c2:00:00:0e). " +
					"STP (bit 0) and MAC pause (bit 1) are always filtered by the kernel.",
			},
		},
	}
}

// bridgeGroupFwdMask returns the configured group_fwd_mask and whether it is
// set for a bridge link.
func bridgeGroupFwdMask(data *models.LinkModel) (uint16, bool) {
	if data.Type.ValueString() != "bridge" || data.Bridge == nil {
		return 0, false
	}
	m := data.Bridge.GroupFwdMask
	if m.IsNull() || m.IsUnknown() {
		return 0, false
	}
	return uint16(m.ValueInt64()), true
}

func buildBridge(name string, cfg *models.BridgeConfig) (*vnl.Bridge, error) {
	br := &vnl.Bridge{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.HelloTime.IsNull() {
			v := uint32(cfg.HelloTime.ValueInt64())
			br.HelloTime = &v
		}
		if !cfg.VlanFiltering.IsNull() {
			v := cfg.VlanFiltering.ValueBool()
			br.VlanFiltering = &v
		}
		if !cfg.AgeingTime.IsNull() {
			v := uint32(cfg.AgeingTime.ValueInt64())
			br.AgeingTime = &v
		}
	}
	return br, nil
}
