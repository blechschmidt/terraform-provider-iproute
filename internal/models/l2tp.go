package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type L2tpTunnelModel struct {
	ID           types.String `tfsdk:"id"`
	TunnelID     types.Int64  `tfsdk:"tunnel_id"`
	PeerTunnelID types.Int64  `tfsdk:"peer_tunnel_id"`
	EncapType    types.String `tfsdk:"encap_type"`
	Local        types.String `tfsdk:"local"`
	Remote       types.String `tfsdk:"remote"`
	LocalPort    types.Int64  `tfsdk:"local_port"`
	RemotePort   types.Int64  `tfsdk:"remote_port"`
}

type L2tpSessionModel struct {
	ID            types.String `tfsdk:"id"`
	TunnelID      types.Int64  `tfsdk:"tunnel_id"`
	SessionID     types.Int64  `tfsdk:"session_id"`
	PeerSessionID types.Int64  `tfsdk:"peer_session_id"`
	Name          types.String `tfsdk:"name"`
	Cookie        types.String `tfsdk:"cookie"`
	PeerCookie    types.String `tfsdk:"peer_cookie"`
	L2specType    types.String `tfsdk:"l2spec_type"`
}
