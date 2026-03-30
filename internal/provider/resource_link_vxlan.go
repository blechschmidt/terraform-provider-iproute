package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func vxlanSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "VxLAN-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"vni":      schema.Int64Attribute{Optional: true, Description: "VxLAN Network Identifier."},
			"group":    schema.StringAttribute{Optional: true, Description: "Multicast group address."},
			"local":    schema.StringAttribute{Optional: true, Description: "Local IP address."},
			"dev":      schema.StringAttribute{Optional: true, Description: "Physical device for tunnel endpoints."},
			"port":     schema.Int64Attribute{Optional: true, Description: "UDP destination port."},
			"learning": schema.BoolAttribute{Optional: true, Description: "Enable MAC learning."},
			"proxy":    schema.BoolAttribute{Optional: true, Description: "Enable ARP proxy."},
			"l2miss":   schema.BoolAttribute{Optional: true, Description: "Enable netlink LLADDR miss notifications."},
			"l3miss":   schema.BoolAttribute{Optional: true, Description: "Enable netlink IP addr miss notifications."},
		},
	}
}

func buildVxlan(name string, cfg *models.VxlanConfig) (*vnl.Vxlan, error) {
	vxlan := &vnl.Vxlan{
		LinkAttrs: vnl.LinkAttrs{Name: name},
	}
	if cfg != nil {
		if !cfg.VNI.IsNull() {
			vxlan.VxlanId = int(cfg.VNI.ValueInt64())
		}
		if !cfg.Group.IsNull() {
			vxlan.Group = net.ParseIP(cfg.Group.ValueString())
		}
		if !cfg.Local.IsNull() {
			vxlan.SrcAddr = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Port.IsNull() {
			vxlan.Port = int(cfg.Port.ValueInt64())
		}
		if !cfg.Learning.IsNull() {
			vxlan.Learning = cfg.Learning.ValueBool()
		}
		if !cfg.Proxy.IsNull() {
			vxlan.Proxy = cfg.Proxy.ValueBool()
		}
		if !cfg.L2miss.IsNull() {
			vxlan.L2miss = cfg.L2miss.ValueBool()
		}
		if !cfg.L3miss.IsNull() {
			vxlan.L3miss = cfg.L3miss.ValueBool()
		}
	}
	return vxlan, nil
}
