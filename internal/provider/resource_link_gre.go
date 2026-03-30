package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func greSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "GRE tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"local":     schema.StringAttribute{Optional: true, Description: "Local endpoint IP."},
			"remote":    schema.StringAttribute{Optional: true, Description: "Remote endpoint IP."},
			"ttl":       schema.Int64Attribute{Optional: true, Description: "TTL."},
			"tos":       schema.Int64Attribute{Optional: true, Description: "TOS."},
			"pmtu_disc": schema.BoolAttribute{Optional: true, Description: "Enable PMTU discovery."},
			"key":       schema.Int64Attribute{Optional: true, Description: "Tunnel key."},
			"ikey":      schema.Int64Attribute{Optional: true, Description: "Input key."},
			"okey":      schema.Int64Attribute{Optional: true, Description: "Output key."},
		},
	}
}

func buildGre(name string, cfg *models.GreConfig) (*vnl.Gretun, error) {
	gre := &vnl.Gretun{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			gre.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			gre.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.TTL.IsNull() {
			gre.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		if !cfg.TOS.IsNull() {
			gre.Tos = uint8(cfg.TOS.ValueInt64())
		}
		if !cfg.PMtuDisc.IsNull() {
			gre.PMtuDisc = 1
		}
		if !cfg.IKey.IsNull() {
			gre.IKey = uint32(cfg.IKey.ValueInt64())
		}
		if !cfg.OKey.IsNull() {
			gre.OKey = uint32(cfg.OKey.ValueInt64())
		}
	}
	return gre, nil
}

func buildGretap(name string, cfg *models.GreConfig) (*vnl.Gretap, error) {
	gre := &vnl.Gretap{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			gre.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			gre.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.TTL.IsNull() {
			gre.Ttl = uint8(cfg.TTL.ValueInt64())
		}
		if !cfg.TOS.IsNull() {
			gre.Tos = uint8(cfg.TOS.ValueInt64())
		}
		if !cfg.PMtuDisc.IsNull() {
			gre.PMtuDisc = 1
		}
		if !cfg.IKey.IsNull() {
			gre.IKey = uint32(cfg.IKey.ValueInt64())
		}
		if !cfg.OKey.IsNull() {
			gre.OKey = uint32(cfg.OKey.ValueInt64())
		}
	}
	return gre, nil
}
