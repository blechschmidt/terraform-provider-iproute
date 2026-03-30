package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func ipipSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "IPIP tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"local":  schema.StringAttribute{Optional: true, Description: "Local endpoint IP."},
			"remote": schema.StringAttribute{Optional: true, Description: "Remote endpoint IP."},
			"ttl":    schema.Int64Attribute{Optional: true, Description: "TTL."},
			"tos":    schema.Int64Attribute{Optional: true, Description: "TOS."},
		},
	}
}

func buildIpip(name string, cfg *models.IpipConfig) (*vnl.Iptun, error) {
	ipip := &vnl.Iptun{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			ipip.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			ipip.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.TTL.IsNull() {
			ipip.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		if !cfg.TOS.IsNull() {
			ipip.Tos = uint8(cfg.TOS.ValueInt64())
		}
	}
	return ipip, nil
}
