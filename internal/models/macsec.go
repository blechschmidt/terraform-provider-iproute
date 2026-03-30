package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MacsecModel struct {
	ID            types.String `tfsdk:"id"`
	Parent        types.String `tfsdk:"parent"`
	Name          types.String `tfsdk:"name"`
	SCI           types.String `tfsdk:"sci"`
	Port          types.Int64  `tfsdk:"port"`
	Encrypt       types.Bool   `tfsdk:"encrypt"`
	CipherSuite   types.String `tfsdk:"cipher_suite"`
	ICVLen        types.Int64  `tfsdk:"icv_len"`
	EncodingSA    types.Int64  `tfsdk:"encoding_sa"`
	Validate      types.String `tfsdk:"validate"`
	Protect       types.Bool   `tfsdk:"protect"`
	ReplayProtect types.Bool   `tfsdk:"replay_protect"`
	Window        types.Int64  `tfsdk:"window"`
}
