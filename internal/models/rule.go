package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RuleModel struct {
	ID                types.String `tfsdk:"id"`
	Priority          types.Int64  `tfsdk:"priority"`
	Family            types.String `tfsdk:"family"`
	Src               types.String `tfsdk:"src"`
	Dst               types.String `tfsdk:"dst"`
	IifName           types.String `tfsdk:"iif_name"`
	OifName           types.String `tfsdk:"oif_name"`
	Table             types.Int64  `tfsdk:"table"`
	FwMark            types.Int64  `tfsdk:"fwmark"`
	FwMask            types.Int64  `tfsdk:"fwmask"`
	TosMatch          types.Int64  `tfsdk:"tos"`
	Action            types.String `tfsdk:"action"`
	Goto              types.Int64  `tfsdk:"goto_priority"`
	SuppressPrefixLen types.Int64  `tfsdk:"suppress_prefix_len"`
	SuppressIfGroup   types.Int64  `tfsdk:"suppress_if_group"`
	IPProto           types.Int64  `tfsdk:"ip_proto"`
	SportRange        types.String `tfsdk:"sport_range"`
	DportRange        types.String `tfsdk:"dport_range"`
	UidRange          types.String `tfsdk:"uid_range"`
	Invert            types.Bool   `tfsdk:"invert"`
}
