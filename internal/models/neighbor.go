package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NeighborModel struct {
	ID       types.String `tfsdk:"id"`
	Address  types.String `tfsdk:"address"`
	LLAddr   types.String `tfsdk:"lladdr"`
	Device   types.String `tfsdk:"device"`
	State    types.String `tfsdk:"state"`
	Family   types.String `tfsdk:"family"`
	Flags    types.List   `tfsdk:"flags"`
	IsRouter types.Bool   `tfsdk:"is_router"`
	Origin   types.String `tfsdk:"origin"`
	Proxy    types.Bool   `tfsdk:"proxy"`
	VNI      types.Int64  `tfsdk:"vni"`
}
