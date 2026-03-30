package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type XfrmStateModel struct {
	ID           types.String  `tfsdk:"id"`
	Src          types.String  `tfsdk:"src"`
	Dst          types.String  `tfsdk:"dst"`
	Proto        types.String  `tfsdk:"proto"`
	SPI          types.Int64   `tfsdk:"spi"`
	Mode         types.String  `tfsdk:"mode"`
	Reqid        types.Int64   `tfsdk:"reqid"`
	ReplayWindow types.Int64   `tfsdk:"replay_window"`
	Auth         *XfrmAlgo     `tfsdk:"auth"`
	Crypt        *XfrmAlgo     `tfsdk:"crypt"`
	Aead         *XfrmAlgoAead `tfsdk:"aead"`
	Mark         types.Int64   `tfsdk:"mark"`
	MarkMask     types.Int64   `tfsdk:"mark_mask"`
	IfID         types.Int64   `tfsdk:"if_id"`
	Family       types.String  `tfsdk:"family"`
}

type XfrmAlgo struct {
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

type XfrmAlgoAead struct {
	Name   types.String `tfsdk:"name"`
	Key    types.String `tfsdk:"key"`
	ICVLen types.Int64  `tfsdk:"icv_len"`
}

type XfrmPolicyModel struct {
	ID       types.String `tfsdk:"id"`
	Src      types.String `tfsdk:"src"`
	Dst      types.String `tfsdk:"dst"`
	Dir      types.String `tfsdk:"dir"`
	Priority types.Int64  `tfsdk:"priority"`
	Action   types.String `tfsdk:"action"`
	Proto    types.String `tfsdk:"proto"`
	SrcPort  types.Int64  `tfsdk:"src_port"`
	DstPort  types.Int64  `tfsdk:"dst_port"`
	Mark     types.Int64  `tfsdk:"mark"`
	MarkMask types.Int64  `tfsdk:"mark_mask"`
	IfID     types.Int64  `tfsdk:"if_id"`
	Family   types.String `tfsdk:"family"`
	Templates types.List  `tfsdk:"templates"`
}

type XfrmPolicyTemplate struct {
	Src   types.String `tfsdk:"src"`
	Dst   types.String `tfsdk:"dst"`
	Proto types.String `tfsdk:"proto"`
	Mode  types.String `tfsdk:"mode"`
	Reqid types.Int64  `tfsdk:"reqid"`
	SPI   types.Int64  `tfsdk:"spi"`
}
