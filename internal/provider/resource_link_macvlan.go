package provider

import (
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func macvlanSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Macvlan-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"parent": schema.StringAttribute{Optional: true, Description: "Parent interface name."},
			"mode":   schema.StringAttribute{Optional: true, Description: "Mode: private, vepa, bridge, passthru, source."},
		},
	}
}

func macvtapSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Macvtap-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"parent": schema.StringAttribute{Optional: true, Description: "Parent interface name."},
			"mode":   schema.StringAttribute{Optional: true, Description: "Mode: private, vepa, bridge, passthru, source."},
		},
	}
}

var macvlanModeMap = map[string]vnl.MacvlanMode{
	"private":  vnl.MACVLAN_MODE_PRIVATE,
	"vepa":     vnl.MACVLAN_MODE_VEPA,
	"bridge":   vnl.MACVLAN_MODE_BRIDGE,
	"passthru": vnl.MACVLAN_MODE_PASSTHRU,
	"source":   vnl.MACVLAN_MODE_SOURCE,
}

func buildMacvlan(name string, cfg *models.MacvlanConfig, client *netlinkClient.Client) (*vnl.Macvlan, error) {
	if cfg == nil {
		return nil, fmt.Errorf("macvlan block is required for macvlan type")
	}
	parent, err := client.LinkByName(cfg.Parent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("parent interface %q not found: %w", cfg.Parent.ValueString(), err)
	}
	mv := &vnl.Macvlan{
		LinkAttrs: vnl.LinkAttrs{Name: name, ParentIndex: parent.Attrs().Index},
		Mode:      vnl.MACVLAN_MODE_BRIDGE,
	}
	if !cfg.Mode.IsNull() {
		if m, ok := macvlanModeMap[cfg.Mode.ValueString()]; ok {
			mv.Mode = m
		}
	}
	return mv, nil
}

func buildMacvtap(name string, cfg *models.MacvtapConfig, client *netlinkClient.Client) (*vnl.Macvtap, error) {
	if cfg == nil {
		return nil, fmt.Errorf("macvtap block is required for macvtap type")
	}
	parent, err := client.LinkByName(cfg.Parent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("parent interface %q not found: %w", cfg.Parent.ValueString(), err)
	}
	mv := &vnl.Macvtap{
		Macvlan: vnl.Macvlan{
			LinkAttrs: vnl.LinkAttrs{Name: name, ParentIndex: parent.Attrs().Index},
			Mode:      vnl.MACVLAN_MODE_BRIDGE,
		},
	}
	if !cfg.Mode.IsNull() {
		if m, ok := macvlanModeMap[cfg.Mode.ValueString()]; ok {
			mv.Macvlan.Mode = m
		}
	}
	return mv, nil
}
