package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func geneveSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Geneve tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"vni":    schema.Int64Attribute{Optional: true, Description: "Virtual Network Identifier."},
			"remote": schema.StringAttribute{Optional: true, Description: "Remote IP address."},
			"port":   schema.Int64Attribute{Optional: true, Description: "UDP port."},
			"ttl":    schema.Int64Attribute{Optional: true, Description: "TTL."},
			"tos":    schema.Int64Attribute{Optional: true, Description: "TOS."},
		},
	}
}

func buildGeneve(name string, cfg *models.GeneveConfig) (*vnl.Geneve, error) {
	gen := &vnl.Geneve{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.VNI.IsNull() {
			gen.ID = uint32(cfg.VNI.ValueInt64())
		}
		if !cfg.Remote.IsNull() {
			gen.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.Port.IsNull() {
			gen.Dport = uint16(cfg.Port.ValueInt64())
		}
		if !cfg.TTL.IsNull() {
			gen.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		if !cfg.TOS.IsNull() {
			gen.Tos = uint8(cfg.TOS.ValueInt64())
		}
	}
	return gen, nil
}
