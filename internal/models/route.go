package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RouteModel struct {
	ID          types.String     `tfsdk:"id"`
	Destination types.String     `tfsdk:"destination"`
	Gateway     types.String     `tfsdk:"gateway"`
	Device      types.String     `tfsdk:"device"`
	Source      types.String     `tfsdk:"source"`
	Metric      types.Int64      `tfsdk:"metric"`
	Table       types.Int64      `tfsdk:"table"`
	Scope       types.String     `tfsdk:"scope"`
	Protocol    types.String     `tfsdk:"protocol"`
	Type        types.String     `tfsdk:"type"`
	MTU         types.Int64      `tfsdk:"mtu"`
	AdvMSS      types.Int64      `tfsdk:"advmss"`
	Family      types.String     `tfsdk:"family"`
	NexthopID   types.Int64      `tfsdk:"nexthop_id"`
	Multipath   types.List       `tfsdk:"multipath"`
	Encap       *RouteEncapModel `tfsdk:"encap"`
	Onlink      types.Bool       `tfsdk:"onlink"`
}

type MultipathModel struct {
	Gateway types.String `tfsdk:"gateway"`
	Device  types.String `tfsdk:"device"`
	Weight  types.Int64  `tfsdk:"weight"`
}

type RouteEncapModel struct {
	Type   types.String `tfsdk:"type"`
	Labels types.List   `tfsdk:"labels"`
}
