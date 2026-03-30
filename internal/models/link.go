package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinkModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	MacAddress  types.String `tfsdk:"mac_address"`
	MTU         types.Int64  `tfsdk:"mtu"`
	TxQueueLen  types.Int64  `tfsdk:"tx_queue_len"`
	Master      types.String `tfsdk:"master"`
	AdminStatus types.String `tfsdk:"admin_status"`
	OperStatus  types.String `tfsdk:"oper_status"`
	IfIndex     types.Int64  `tfsdk:"if_index"`
	Speed       types.Int64  `tfsdk:"speed"`
	Bridge      *BridgeConfig    `tfsdk:"bridge"`
	Veth        *VethConfig      `tfsdk:"veth"`
	Vlan        *VlanConfig      `tfsdk:"vlan"`
	Vxlan       *VxlanConfig     `tfsdk:"vxlan"`
	Bond        *BondConfig      `tfsdk:"bond"`
	Macvlan     *MacvlanConfig   `tfsdk:"macvlan"`
	Macvtap     *MacvtapConfig   `tfsdk:"macvtap"`
	Ipvlan      *IpvlanConfig    `tfsdk:"ipvlan"`
	Dummy       *DummyConfig     `tfsdk:"dummy"`
	Ifb         *IfbConfig       `tfsdk:"ifb"`
	Gre         *GreConfig       `tfsdk:"gre"`
	Sit         *SitConfig       `tfsdk:"sit"`
	Vti         *VtiConfig       `tfsdk:"vti"`
	Ip6tnl      *Ip6tnlConfig    `tfsdk:"ip6tnl"`
	Ipip        *IpipConfig      `tfsdk:"ipip"`
	Geneve      *GeneveConfig    `tfsdk:"geneve"`
	Wireguard   *WireguardConfig `tfsdk:"wireguard"`
	Tuntap      *TuntapLinkConfig `tfsdk:"tuntap"`
	Statistics  types.Object     `tfsdk:"statistics"`
}

type BridgeConfig struct {
	STP          types.Bool  `tfsdk:"stp"`
	HelloTime    types.Int64 `tfsdk:"hello_time"`
	MaxAge       types.Int64 `tfsdk:"max_age"`
	ForwardDelay types.Int64 `tfsdk:"forward_delay"`
	VlanFiltering types.Bool `tfsdk:"vlan_filtering"`
	DefaultPVID  types.Int64 `tfsdk:"default_pvid"`
	AgeingTime   types.Int64 `tfsdk:"ageing_time"`
}

type VethConfig struct {
	PeerName types.String `tfsdk:"peer_name"`
}

type VlanConfig struct {
	VlanID   types.Int64  `tfsdk:"vlan_id"`
	Parent   types.String `tfsdk:"parent"`
	Protocol types.String `tfsdk:"protocol"`
}

type VxlanConfig struct {
	VNI      types.Int64  `tfsdk:"vni"`
	Group    types.String `tfsdk:"group"`
	Local    types.String `tfsdk:"local"`
	Dev      types.String `tfsdk:"dev"`
	Port     types.Int64  `tfsdk:"port"`
	Learning types.Bool   `tfsdk:"learning"`
	Proxy    types.Bool   `tfsdk:"proxy"`
	L2miss   types.Bool   `tfsdk:"l2miss"`
	L3miss   types.Bool   `tfsdk:"l3miss"`
}

type BondConfig struct {
	Mode           types.String `tfsdk:"mode"`
	MiiMon         types.Int64  `tfsdk:"miimon"`
	UpDelay        types.Int64  `tfsdk:"up_delay"`
	DownDelay      types.Int64  `tfsdk:"down_delay"`
	Primary        types.String `tfsdk:"primary"`
	LacpRate       types.String `tfsdk:"lacp_rate"`
	XmitHashPolicy types.String `tfsdk:"xmit_hash_policy"`
}

type MacvlanConfig struct {
	Parent types.String `tfsdk:"parent"`
	Mode   types.String `tfsdk:"mode"`
}

type MacvtapConfig struct {
	Parent types.String `tfsdk:"parent"`
	Mode   types.String `tfsdk:"mode"`
}

type IpvlanConfig struct {
	Parent types.String `tfsdk:"parent"`
	Mode   types.String `tfsdk:"mode"`
	Flag   types.String `tfsdk:"flag"`
}

type DummyConfig struct{}

type IfbConfig struct{}

type GreConfig struct {
	Local    types.String `tfsdk:"local"`
	Remote   types.String `tfsdk:"remote"`
	TTL      types.Int64  `tfsdk:"ttl"`
	TOS      types.Int64  `tfsdk:"tos"`
	PMtuDisc types.Bool   `tfsdk:"pmtu_disc"`
	Key      types.Int64  `tfsdk:"key"`
	IKey     types.Int64  `tfsdk:"ikey"`
	OKey     types.Int64  `tfsdk:"okey"`
}

type SitConfig struct {
	Local  types.String `tfsdk:"local"`
	Remote types.String `tfsdk:"remote"`
	TTL    types.Int64  `tfsdk:"ttl"`
	TOS    types.Int64  `tfsdk:"tos"`
}

type VtiConfig struct {
	Local  types.String `tfsdk:"local"`
	Remote types.String `tfsdk:"remote"`
	IKey   types.Int64  `tfsdk:"ikey"`
	OKey   types.Int64  `tfsdk:"okey"`
}

type Ip6tnlConfig struct {
	Local      types.String `tfsdk:"local"`
	Remote     types.String `tfsdk:"remote"`
	TTL        types.Int64  `tfsdk:"ttl"`
	FlowLabel  types.Int64  `tfsdk:"flow_label"`
	EncapLimit types.Int64  `tfsdk:"encap_limit"`
}

type IpipConfig struct {
	Local  types.String `tfsdk:"local"`
	Remote types.String `tfsdk:"remote"`
	TTL    types.Int64  `tfsdk:"ttl"`
	TOS    types.Int64  `tfsdk:"tos"`
}

type GeneveConfig struct {
	VNI    types.Int64  `tfsdk:"vni"`
	Remote types.String `tfsdk:"remote"`
	Port   types.Int64  `tfsdk:"port"`
	TTL    types.Int64  `tfsdk:"ttl"`
	TOS    types.Int64  `tfsdk:"tos"`
}

type WireguardConfig struct {
	PrivateKey types.String `tfsdk:"private_key"`
	ListenPort types.Int64  `tfsdk:"listen_port"`
	FwMark     types.Int64  `tfsdk:"fwmark"`
	Peers      types.List   `tfsdk:"peers"`
}

type WireguardPeerModel struct {
	PublicKey           types.String `tfsdk:"public_key"`
	PresharedKey        types.String `tfsdk:"preshared_key"`
	Endpoint            types.String `tfsdk:"endpoint"`
	AllowedIPs          types.List   `tfsdk:"allowed_ips"`
	PersistentKeepalive types.Int64  `tfsdk:"persistent_keepalive"`
}

type TuntapLinkConfig struct {
	Mode       types.String `tfsdk:"mode"`
	Owner      types.Int64  `tfsdk:"owner"`
	Group      types.Int64  `tfsdk:"group"`
	MultiQueue types.Bool   `tfsdk:"multi_queue"`
}

type LinkStatistics struct {
	RxBytes   types.Int64 `tfsdk:"rx_bytes"`
	TxBytes   types.Int64 `tfsdk:"tx_bytes"`
	RxPackets types.Int64 `tfsdk:"rx_packets"`
	TxPackets types.Int64 `tfsdk:"tx_packets"`
	RxErrors  types.Int64 `tfsdk:"rx_errors"`
	TxErrors  types.Int64 `tfsdk:"tx_errors"`
	RxDropped types.Int64 `tfsdk:"rx_dropped"`
	TxDropped types.Int64 `tfsdk:"tx_dropped"`
	Multicast types.Int64 `tfsdk:"multicast"`
}
