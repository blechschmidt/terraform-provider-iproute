package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TcpMetricsModel struct {
	ID       types.String `tfsdk:"id"`
	Address  types.String `tfsdk:"address"`
	RTT      types.Int64  `tfsdk:"rtt"`
	RTTVar   types.Int64  `tfsdk:"rttvar"`
	SSTHRESH types.Int64  `tfsdk:"ssthresh"`
	CWND     types.Int64  `tfsdk:"cwnd"`
}
