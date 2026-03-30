package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TokenModel struct {
	ID     types.String `tfsdk:"id"`
	Device types.String `tfsdk:"device"`
	Token  types.String `tfsdk:"token"`
}
