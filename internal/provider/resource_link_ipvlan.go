package provider

import (
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func ipvlanSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "IPVlan-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"parent": schema.StringAttribute{Optional: true, Description: "Parent interface name."},
			"mode":   schema.StringAttribute{Optional: true, Description: "Mode: l2, l3, l3s."},
			"flag":   schema.StringAttribute{Optional: true, Description: "Flag: bridge, private, vepa."},
		},
	}
}

func buildIpvlan(name string, cfg *models.IpvlanConfig, client *netlinkClient.Client) (*vnl.IPVlan, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ipvlan block is required for ipvlan type")
	}
	parent, err := client.LinkByName(cfg.Parent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("parent interface %q not found: %w", cfg.Parent.ValueString(), err)
	}
	ipv := &vnl.IPVlan{
		LinkAttrs: vnl.LinkAttrs{Name: name, ParentIndex: parent.Attrs().Index},
		Mode:      vnl.IPVLAN_MODE_L2,
	}
	if !cfg.Mode.IsNull() {
		switch cfg.Mode.ValueString() {
		case "l3":
			ipv.Mode = vnl.IPVLAN_MODE_L3
		case "l3s":
			ipv.Mode = vnl.IPVLAN_MODE_L3S
		}
	}
	if !cfg.Flag.IsNull() {
		switch cfg.Flag.ValueString() {
		case "bridge":
			ipv.Flag = vnl.IPVLAN_FLAG_BRIDGE
		case "private":
			ipv.Flag = vnl.IPVLAN_FLAG_PRIVATE
		case "vepa":
			ipv.Flag = vnl.IPVLAN_FLAG_VEPA
		}
	}
	return ipv, nil
}
