package provider

import (
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func vlanSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "VLAN-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"vlan_id":  schema.Int64Attribute{Optional: true, Description: "VLAN ID (1-4094)."},
			"parent":   schema.StringAttribute{Optional: true, Description: "Parent interface name."},
			"protocol": schema.StringAttribute{Optional: true, Description: "VLAN protocol (802.1Q or 802.1ad)."},
		},
	}
}

func buildVlan(name string, cfg *models.VlanConfig, client *netlinkClient.Client) (*vnl.Vlan, error) {
	if cfg == nil {
		return nil, fmt.Errorf("vlan block is required for vlan type")
	}
	parent, err := client.LinkByName(cfg.Parent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("parent interface %q not found: %w", cfg.Parent.ValueString(), err)
	}
	vlan := &vnl.Vlan{
		LinkAttrs:    vnl.LinkAttrs{Name: name, ParentIndex: parent.Attrs().Index},
		VlanId:       int(cfg.VlanID.ValueInt64()),
		VlanProtocol: vnl.VLAN_PROTOCOL_8021Q,
	}
	if !cfg.Protocol.IsNull() && cfg.Protocol.ValueString() == "802.1ad" {
		vlan.VlanProtocol = vnl.VLAN_PROTOCOL_8021AD
	}
	return vlan, nil
}
