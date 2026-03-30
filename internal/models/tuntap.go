package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TuntapModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Mode       types.String `tfsdk:"mode"`
	Owner      types.Int64  `tfsdk:"owner"`
	Group      types.Int64  `tfsdk:"group"`
	MultiQueue types.Bool   `tfsdk:"multi_queue"`
	Persist    types.Bool   `tfsdk:"persist"`
}
