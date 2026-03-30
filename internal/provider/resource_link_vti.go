package provider

import (
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func vtiSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "VTI tunnel configuration.",
		Attributes: map[string]schema.Attribute{
			"local":  schema.StringAttribute{Optional: true, Description: "Local endpoint IP."},
			"remote": schema.StringAttribute{Optional: true, Description: "Remote endpoint IP."},
			"ikey":   schema.Int64Attribute{Optional: true, Description: "Input key."},
			"okey":   schema.Int64Attribute{Optional: true, Description: "Output key."},
		},
	}
}

func buildVti(name string, cfg *models.VtiConfig) (*vnl.Vti, error) {
	vti := &vnl.Vti{LinkAttrs: vnl.LinkAttrs{Name: name}}
	if cfg != nil {
		if !cfg.Local.IsNull() {
			vti.Local = net.ParseIP(cfg.Local.ValueString())
		}
		if !cfg.Remote.IsNull() {
			vti.Remote = net.ParseIP(cfg.Remote.ValueString())
		}
		if !cfg.IKey.IsNull() {
			vti.IKey = uint32(cfg.IKey.ValueInt64())
		}
		if !cfg.OKey.IsNull() {
			vti.OKey = uint32(cfg.OKey.ValueInt64())
		}
	}
	return vti, nil
}
