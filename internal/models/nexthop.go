package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NexthopModel struct {
	ID        types.String `tfsdk:"id"`
	NhID      types.Int64  `tfsdk:"nhid"`
	Gateway   types.String `tfsdk:"gateway"`
	Device    types.String `tfsdk:"device"`
	Blackhole types.Bool   `tfsdk:"blackhole"`
	Family    types.String `tfsdk:"family"`
	Group     types.List   `tfsdk:"group"`
	Resilient types.Bool   `tfsdk:"resilient"`
	FDB       types.Bool   `tfsdk:"fdb"`
}

type NexthopGroupMemberModel struct {
	ID     types.Int64 `tfsdk:"id"`
	Weight types.Int64 `tfsdk:"weight"`
}
