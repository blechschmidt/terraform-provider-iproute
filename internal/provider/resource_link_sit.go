package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func sitSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "SIT tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"local":  schema.StringAttribute{Optional: true, Description: "Local endpoint IP."},
			"remote": schema.StringAttribute{Optional: true, Description: "Remote endpoint IP."},
			"ttl":    schema.Int64Attribute{Optional: true, Description: "TTL."},
			"tos":    schema.Int64Attribute{Optional: true, Description: "TOS."},
		},
	}
}

func buildSit(name string, cfg *models.SitConfig) (*vnl.Sittun, error) {
	sit := &vnl.Sittun{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			sit.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			sit.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.TTL.IsNull() {
			sit.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		if !cfg.TOS.IsNull() {
			sit.Tos = uint8(cfg.TOS.ValueInt64())
		}
	}
	return sit, nil
}
