package provider

import (
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func tuntapSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "TUN/TAP configuration.",
		Attributes: map[string]schema.Attribute{
			"mode":        schema.StringAttribute{Optional: true, Description: "Mode: tun or tap."},
			"owner":       schema.Int64Attribute{Optional: true, Description: "Owner UID."},
			"group":       schema.Int64Attribute{Optional: true, Description: "Group GID."},
			"multi_queue": schema.BoolAttribute{Optional: true, Description: "Enable multi-queue."},
		},
	}
}

func buildTuntap(name string, cfg *models.TuntapLinkConfig) (*vnl.Tuntap, error) {
	if cfg == nil {
		return nil, fmt.Errorf("tuntap block is required for tuntap type")
	}
	mode := vnl.TUNTAP_MODE_TUN
	if cfg.Mode.ValueString() == "tap" {
		mode = vnl.TUNTAP_MODE_TAP
	}
	tt := &vnl.Tuntap{
		LinkAttrs: vnl.LinkAttrs{Name: name},
		Mode:      mode,
	}
	if !cfg.Owner.IsNull() {
		v := int(cfg.Owner.ValueInt64())
		tt.Owner = uint32(v)
	}
	if !cfg.Group.IsNull() {
		v := int(cfg.Group.ValueInt64())
		tt.Group = uint32(v)
	}
	if !cfg.MultiQueue.IsNull() && cfg.MultiQueue.ValueBool() {
		tt.Flags = vnl.TUNTAP_MULTI_QUEUE_DEFAULTS
	}
	return tt, nil
}
