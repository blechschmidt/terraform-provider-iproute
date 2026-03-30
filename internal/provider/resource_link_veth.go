package provider

import (
	"github.com/example/terraform-provider-iproute/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	vnl "github.com/vishvananda/netlink"
)

func vethSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Veth-specific configuration.",
		Attributes: map[string]schema.Attribute{
			"peer_name": schema.StringAttribute{Optional: true, Description: "Peer interface name."},
		},
	}
}

func buildVeth(name string, cfg *models.VethConfig) (*vnl.Veth, error) {
	peerName := "veth0"
	if cfg != nil && !cfg.PeerName.IsNull() {
		peerName = cfg.PeerName.ValueString()
	}
	return &vnl.Veth{
		LinkAttrs: vnl.LinkAttrs{Name: name},
		PeerName:  peerName,
	}, nil
}
