package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FouModel struct {
	ID         types.String `tfsdk:"id"`
	Port       types.Int64  `tfsdk:"port"`
	Family     types.String `tfsdk:"family"`
	Protocol   types.Int64  `tfsdk:"protocol"`
	EncapType  types.String `tfsdk:"encap_type"`
	RemotePort types.Int64  `tfsdk:"remote_port"`
	Local      types.String `tfsdk:"local"`
	Peer       types.String `tfsdk:"peer"`
}
