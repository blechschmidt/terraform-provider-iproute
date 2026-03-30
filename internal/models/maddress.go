package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MaddressModel struct {
	ID      types.String `tfsdk:"id"`
	Device  types.String `tfsdk:"device"`
	Address types.String `tfsdk:"address"`
}
