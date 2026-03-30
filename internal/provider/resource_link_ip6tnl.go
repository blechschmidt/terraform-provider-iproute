package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func ip6tnlSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "IP6TNL tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"local":       schema.StringAttribute{Optional: true, Description: "Local IPv6 endpoint."},
			"remote":      schema.StringAttribute{Optional: true, Description: "Remote IPv6 endpoint."},
			"ttl":         schema.Int64Attribute{Optional: true, Description: "TTL/Hop limit."},
			"flow_label":  schema.Int64Attribute{Optional: true, Description: "Flow label."},
			"encap_limit": schema.Int64Attribute{Optional: true, Description: "Encap limit."},
		},
	}
}

func buildIp6tnl(name string, cfg *models.Ip6tnlConfig) (*vnl.Ip6tnl, error) {
	tnl := &vnl.Ip6tnl{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			tnl.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			tnl.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.TTL.IsNull() {
			tnl.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		// FlowLabel not directly supported, skipped
		if !cfg.EncapLimit.IsNull() {
			tnl.EncapLimit = uint8(cfg.EncapLimit.ValueInt64())
		}
	}
	return tnl, nil
}
