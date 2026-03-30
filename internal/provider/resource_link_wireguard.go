package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func wireguardSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "WireGuard configuration. Key/peer management via wgctrl.",
		Attributes: map[string]schema.Attribute{
			"private_key":  schema.StringAttribute{Optional: true, Sensitive: true, Description: "Private key."},
			"listen_port":  schema.Int64Attribute{Optional: true, Description: "Listen port."},
			"fwmark":       schema.Int64Attribute{Optional: true, Description: "Firewall mark."},
			"peers": schema.ListNestedAttribute{
				Optional:    true,
				Description: "WireGuard peers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"public_key":           schema.StringAttribute{Required: true, Description: "Public key."},
						"preshared_key":        schema.StringAttribute{Optional: true, Sensitive: true, Description: "Preshared key."},
						"endpoint":             schema.StringAttribute{Optional: true, Description: "Endpoint (host:port)."},
						"allowed_ips":          schema.ListAttribute{ElementType: types.StringType, Optional: true, Description: "Allowed IP ranges."},
						"persistent_keepalive": schema.Int64Attribute{Optional: true, Description: "Keepalive interval (seconds)."},
					},
				},
			},
		},
	}
}
