package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetnsModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
