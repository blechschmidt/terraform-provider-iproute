package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SrModel struct {
	ID       types.String `tfsdk:"id"`
	Device   types.String `tfsdk:"device"`
	Hmac     types.String `tfsdk:"hmac"`
	Segments types.List   `tfsdk:"segments"`
	Encap    types.String `tfsdk:"encap"`
}
