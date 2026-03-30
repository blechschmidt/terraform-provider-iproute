package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AddressModel struct {
	ID                types.String `tfsdk:"id"`
	Address           types.String `tfsdk:"address"`
	Device            types.String `tfsdk:"device"`
	Peer              types.String `tfsdk:"peer"`
	Broadcast         types.String `tfsdk:"broadcast"`
	Label             types.String `tfsdk:"label"`
	Scope             types.String `tfsdk:"scope"`
	Family            types.String `tfsdk:"family"`
	Origin            types.String `tfsdk:"origin"`
	Flags             types.List   `tfsdk:"flags"`
	PreferredLifetime types.Int64  `tfsdk:"preferred_lifetime"`
	ValidLifetime     types.Int64  `tfsdk:"valid_lifetime"`
}
