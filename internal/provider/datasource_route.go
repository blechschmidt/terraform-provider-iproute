package provider

import (
	"context"
	"fmt"
	"net"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var _ datasource.DataSource = &RouteDataSource{}

type RouteDataSource struct{ client *netlinkClient.Client }

func NewRouteDataSource() datasource.DataSource { return &RouteDataSource{} }

func (d *RouteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

var (
	multipathObjectAttrTypes = map[string]attr.Type{
		"link_index": types.Int64Type,
		"device":     types.StringType,
		"gateway":    types.StringType,
		"hops":       types.Int64Type,
		"weight":     types.Int64Type,
		"flags":      types.ListType{ElemType: types.StringType},
		"new_dst":    types.StringType,
		"via":        types.StringType,
		"encap_type": types.StringType,
		"encap":      types.StringType,
	}

	encapObjectAttrTypes = map[string]attr.Type{
		"type":   types.StringType,
		"string": types.StringType,
	}

	routeObjectAttrTypes = map[string]attr.Type{
		"link_index":          types.Int64Type,
		"device":              types.StringType,
		"ilink_index":         types.Int64Type,
		"in_device":           types.StringType,
		"scope":               types.StringType,
		"destination":         types.StringType,
		"source":              types.StringType,
		"gateway":             types.StringType,
		"protocol":            types.StringType,
		"priority":            types.Int64Type,
		"family":              types.StringType,
		"table":               types.Int64Type,
		"type":                types.StringType,
		"tos":                 types.Int64Type,
		"flags":               types.ListType{ElemType: types.StringType},
		"mpls_dst":            types.Int64Type,
		"new_dst":             types.StringType,
		"via":                 types.StringType,
		"realm":               types.Int64Type,
		"mtu":                 types.Int64Type,
		"window":              types.Int64Type,
		"rtt":                 types.Int64Type,
		"rtt_var":             types.Int64Type,
		"ssthresh":            types.Int64Type,
		"cwnd":                types.Int64Type,
		"advmss":              types.Int64Type,
		"reordering":          types.Int64Type,
		"hoplimit":            types.Int64Type,
		"init_cwnd":           types.Int64Type,
		"features":            types.Int64Type,
		"rto_min":             types.Int64Type,
		"init_rwnd":           types.Int64Type,
		"quick_ack":           types.Int64Type,
		"congctl":             types.StringType,
		"fast_open_no_cookie": types.Int64Type,
		"encap":               types.ObjectType{AttrTypes: encapObjectAttrTypes},
		"multipath":           types.ListType{ElemType: types.ObjectType{AttrTypes: multipathObjectAttrTypes}},
	}
)

func (d *RouteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	multipathAttrs := map[string]schema.Attribute{
		"link_index": schema.Int64Attribute{Computed: true, Description: "Output interface index."},
		"device":     schema.StringAttribute{Computed: true, Description: "Output interface name."},
		"gateway":    schema.StringAttribute{Computed: true, Description: "Gateway IP address."},
		"hops":       schema.Int64Attribute{Computed: true, Description: "Hop count (weight is hops+1)."},
		"weight":     schema.Int64Attribute{Computed: true, Description: "Path weight (hops+1)."},
		"flags":      schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "Nexthop flags."},
		"new_dst":    schema.StringAttribute{Computed: true, Description: "New destination (MPLS swap)."},
		"via":        schema.StringAttribute{Computed: true, Description: "Via address for cross-family nexthops."},
		"encap_type": schema.StringAttribute{Computed: true, Description: "Encapsulation type."},
		"encap":      schema.StringAttribute{Computed: true, Description: "Encapsulation string representation."},
	}

	routeAttrs := map[string]schema.Attribute{
		"link_index":          schema.Int64Attribute{Computed: true, Description: "Output interface index."},
		"device":              schema.StringAttribute{Computed: true, Description: "Output interface name."},
		"ilink_index":         schema.Int64Attribute{Computed: true, Description: "Input interface index (iif)."},
		"in_device":           schema.StringAttribute{Computed: true, Description: "Input interface name (iif)."},
		"scope":               schema.StringAttribute{Computed: true, Description: "Route scope (global, site, link, host, nowhere)."},
		"destination":         schema.StringAttribute{Computed: true, Description: "Destination prefix in CIDR notation, or 'default'."},
		"source":              schema.StringAttribute{Computed: true, Description: "Preferred source address (prefsrc)."},
		"gateway":             schema.StringAttribute{Computed: true, Description: "Gateway IP address."},
		"protocol":            schema.StringAttribute{Computed: true, Description: "Route protocol (kernel, boot, static, ...)."},
		"priority":            schema.Int64Attribute{Computed: true, Description: "Route metric/priority."},
		"family":              schema.StringAttribute{Computed: true, Description: "Address family (inet, inet6, mpls)."},
		"table":               schema.Int64Attribute{Computed: true, Description: "Routing table ID."},
		"type":                schema.StringAttribute{Computed: true, Description: "Route type (unicast, blackhole, ...)."},
		"tos":                 schema.Int64Attribute{Computed: true, Description: "Type of service."},
		"flags":               schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "Route flags."},
		"mpls_dst":            schema.Int64Attribute{Computed: true, Description: "MPLS destination label."},
		"new_dst":             schema.StringAttribute{Computed: true, Description: "New destination (MPLS swap)."},
		"via":                 schema.StringAttribute{Computed: true, Description: "Via address for cross-family nexthops."},
		"realm":               schema.Int64Attribute{Computed: true, Description: "Routing realm."},
		"mtu":                 schema.Int64Attribute{Computed: true, Description: "Route MTU metric."},
		"window":              schema.Int64Attribute{Computed: true, Description: "Advertised window."},
		"rtt":                 schema.Int64Attribute{Computed: true, Description: "Initial round-trip time."},
		"rtt_var":             schema.Int64Attribute{Computed: true, Description: "Initial round-trip time variance."},
		"ssthresh":            schema.Int64Attribute{Computed: true, Description: "Slow-start threshold."},
		"cwnd":                schema.Int64Attribute{Computed: true, Description: "Congestion window."},
		"advmss":              schema.Int64Attribute{Computed: true, Description: "Advertised MSS."},
		"reordering":          schema.Int64Attribute{Computed: true, Description: "Maximum reordering."},
		"hoplimit":            schema.Int64Attribute{Computed: true, Description: "Hop limit (IPv6) / TTL (IPv4)."},
		"init_cwnd":           schema.Int64Attribute{Computed: true, Description: "Initial congestion window."},
		"features":            schema.Int64Attribute{Computed: true, Description: "Route feature flags."},
		"rto_min":             schema.Int64Attribute{Computed: true, Description: "Minimum retransmit timeout."},
		"init_rwnd":           schema.Int64Attribute{Computed: true, Description: "Initial receive window."},
		"quick_ack":           schema.Int64Attribute{Computed: true, Description: "Quick ACK enable."},
		"congctl":             schema.StringAttribute{Computed: true, Description: "Congestion control algorithm."},
		"fast_open_no_cookie": schema.Int64Attribute{Computed: true, Description: "TCP Fast Open without cookie."},
		"encap": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Route encapsulation.",
			Attributes: map[string]schema.Attribute{
				"type":   schema.StringAttribute{Computed: true, Description: "Encapsulation type (mpls, ip, ip6, ila, bpf, seg6, seg6local, xfrm)."},
				"string": schema.StringAttribute{Computed: true, Description: "String representation of the encapsulation."},
			},
		},
		"multipath": schema.ListNestedAttribute{
			Computed:     true,
			Description:  "Multipath nexthops (ECMP).",
			NestedObject: schema.NestedAttributeObject{Attributes: multipathAttrs},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Read routing table entries. Returns all netlink attributes for each route. " +
			"When `get` is set, the data source performs a route lookup for the given destination " +
			"(equivalent to `ip route get <addr>`).",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"family": schema.StringAttribute{Optional: true, Description: "Address family filter (inet, inet6). Ignored when `get` is set; the family is then derived from the address."},
			"table":  schema.Int64Attribute{Optional: true, Description: "Filter by routing table ID. Ignored when `get` is set."},
			"get":    schema.StringAttribute{Optional: true, Description: "Resolve the route used to reach the given IP address (like `ip route get`)."},
			"routes": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of routes with full netlink attributes.",
				NestedObject: schema.NestedAttributeObject{Attributes: routeAttrs},
			},
		},
	}
}

func (d *RouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netlinkClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected type", fmt.Sprintf("Expected *netlink.Client, got %T", req.ProviderData))
		return
	}
	d.client = client
}

type routeDS struct {
	ID     types.String `tfsdk:"id"`
	Family types.String `tfsdk:"family"`
	Table  types.Int64  `tfsdk:"table"`
	Get    types.String `tfsdk:"get"`
	Routes types.List   `tfsdk:"routes"`
}

func (d *RouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data routeDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var routes []vnl.Route
	var err error

	if !data.Get.IsNull() && !data.Get.IsUnknown() && data.Get.ValueString() != "" {
		dst := net.ParseIP(data.Get.ValueString())
		if dst == nil {
			resp.Diagnostics.AddError("Invalid get address", fmt.Sprintf("%q is not a valid IP address", data.Get.ValueString()))
			return
		}
		routes, err = d.client.RouteGet(dst)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get route", err.Error())
			return
		}
		data.ID = types.StringValue("get:" + dst.String())
	} else {
		family := unix.AF_UNSPEC
		if !data.Family.IsNull() && !data.Family.IsUnknown() {
			switch data.Family.ValueString() {
			case "inet":
				family = unix.AF_INET
			case "inet6":
				family = unix.AF_INET6
			}
		}

		filter := &vnl.Route{}
		var mask uint64
		if !data.Table.IsNull() && !data.Table.IsUnknown() {
			filter.Table = int(data.Table.ValueInt64())
			mask |= vnl.RT_FILTER_TABLE
		}

		if mask != 0 {
			routes, err = d.client.RouteListFiltered(family, filter, mask)
		} else {
			routes, err = d.client.RouteList(nil, family)
		}
		if err != nil {
			resp.Diagnostics.AddError("Failed to list routes", err.Error())
			return
		}
		data.ID = types.StringValue("routes")
	}

	values := make([]attr.Value, 0, len(routes))
	for i := range routes {
		obj, diags := d.routeToObject(ctx, &routes[i])
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	list, diags := types.ListValue(types.ObjectType{AttrTypes: routeObjectAttrTypes}, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Routes = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RouteDataSource) routeToObject(ctx context.Context, r *vnl.Route) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	destination := "default"
	if r.Dst != nil {
		destination = r.Dst.String()
	}
	source := ""
	if r.Src != nil {
		source = r.Src.String()
	}
	gateway := ""
	if r.Gw != nil {
		gateway = r.Gw.String()
	}

	device := ""
	if r.LinkIndex > 0 {
		if link, err := d.client.LinkByIndex(r.LinkIndex); err == nil {
			device = link.Attrs().Name
		}
	}
	inDevice := ""
	if r.ILinkIndex > 0 {
		if link, err := d.client.LinkByIndex(r.ILinkIndex); err == nil {
			inDevice = link.Attrs().Name
		}
	}

	family := familyToString(r.Family)
	if family == "" {
		if r.Dst != nil {
			if r.Dst.IP.To4() != nil {
				family = "inet"
			} else {
				family = "inet6"
			}
		} else if r.Gw != nil {
			if r.Gw.To4() != nil {
				family = "inet"
			} else {
				family = "inet6"
			}
		}
	}

	flagsList, fd := stringListValue(ctx, r.ListFlags())
	diags.Append(fd...)

	mplsDst := types.Int64Null()
	if r.MPLSDst != nil {
		mplsDst = types.Int64Value(int64(*r.MPLSDst))
	}

	newDst := ""
	if r.NewDst != nil {
		newDst = r.NewDst.String()
	}
	via := ""
	if r.Via != nil {
		via = r.Via.String()
	}

	encap := types.ObjectNull(encapObjectAttrTypes)
	if r.Encap != nil {
		encapVal, ed := types.ObjectValue(encapObjectAttrTypes, map[string]attr.Value{
			"type":   types.StringValue(encapTypeToString(r.Encap.Type())),
			"string": types.StringValue(r.Encap.String()),
		})
		diags.Append(ed...)
		encap = encapVal
	}

	mpVals := make([]attr.Value, 0, len(r.MultiPath))
	for _, nh := range r.MultiPath {
		mpObj, md := nexthopInfoToObject(ctx, d.client, nh)
		diags.Append(md...)
		mpVals = append(mpVals, mpObj)
	}
	multipath, md := types.ListValue(types.ObjectType{AttrTypes: multipathObjectAttrTypes}, mpVals)
	diags.Append(md...)

	values := map[string]attr.Value{
		"link_index":          types.Int64Value(int64(r.LinkIndex)),
		"device":              types.StringValue(device),
		"ilink_index":         types.Int64Value(int64(r.ILinkIndex)),
		"in_device":           types.StringValue(inDevice),
		"scope":               types.StringValue(scopeToString(int(r.Scope))),
		"destination":         types.StringValue(destination),
		"source":              types.StringValue(source),
		"gateway":             types.StringValue(gateway),
		"protocol":            types.StringValue(protocolToString(r.Protocol)),
		"priority":            types.Int64Value(int64(r.Priority)),
		"family":              types.StringValue(family),
		"table":               types.Int64Value(int64(r.Table)),
		"type":                types.StringValue(routeTypeToString(r.Type)),
		"tos":                 types.Int64Value(int64(r.Tos)),
		"flags":               flagsList,
		"mpls_dst":            mplsDst,
		"new_dst":             types.StringValue(newDst),
		"via":                 types.StringValue(via),
		"realm":               types.Int64Value(int64(r.Realm)),
		"mtu":                 types.Int64Value(int64(r.MTU)),
		"window":              types.Int64Value(int64(r.Window)),
		"rtt":                 types.Int64Value(int64(r.Rtt)),
		"rtt_var":             types.Int64Value(int64(r.RttVar)),
		"ssthresh":            types.Int64Value(int64(r.Ssthresh)),
		"cwnd":                types.Int64Value(int64(r.Cwnd)),
		"advmss":              types.Int64Value(int64(r.AdvMSS)),
		"reordering":          types.Int64Value(int64(r.Reordering)),
		"hoplimit":            types.Int64Value(int64(r.Hoplimit)),
		"init_cwnd":           types.Int64Value(int64(r.InitCwnd)),
		"features":            types.Int64Value(int64(r.Features)),
		"rto_min":             types.Int64Value(int64(r.RtoMin)),
		"init_rwnd":           types.Int64Value(int64(r.InitRwnd)),
		"quick_ack":           types.Int64Value(int64(r.QuickACK)),
		"congctl":             types.StringValue(r.Congctl),
		"fast_open_no_cookie": types.Int64Value(int64(r.FastOpenNoCookie)),
		"encap":               encap,
		"multipath":           multipath,
	}

	obj, od := types.ObjectValue(routeObjectAttrTypes, values)
	diags.Append(od...)
	return obj, diags
}

func nexthopInfoToObject(ctx context.Context, client *netlinkClient.Client, nh *vnl.NexthopInfo) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	device := ""
	if nh.LinkIndex > 0 {
		if link, err := client.LinkByIndex(nh.LinkIndex); err == nil {
			device = link.Attrs().Name
		}
	}
	gateway := ""
	if nh.Gw != nil {
		gateway = nh.Gw.String()
	}
	newDst := ""
	if nh.NewDst != nil {
		newDst = nh.NewDst.String()
	}
	via := ""
	if nh.Via != nil {
		via = nh.Via.String()
	}
	encapType := ""
	encapStr := ""
	if nh.Encap != nil {
		encapType = encapTypeToString(nh.Encap.Type())
		encapStr = nh.Encap.String()
	}

	flagsList, fd := stringListValue(ctx, nh.ListFlags())
	diags.Append(fd...)

	values := map[string]attr.Value{
		"link_index": types.Int64Value(int64(nh.LinkIndex)),
		"device":     types.StringValue(device),
		"gateway":    types.StringValue(gateway),
		"hops":       types.Int64Value(int64(nh.Hops)),
		"weight":     types.Int64Value(int64(nh.Hops + 1)),
		"flags":      flagsList,
		"new_dst":    types.StringValue(newDst),
		"via":        types.StringValue(via),
		"encap_type": types.StringValue(encapType),
		"encap":      types.StringValue(encapStr),
	}

	obj, od := types.ObjectValue(multipathObjectAttrTypes, values)
	diags.Append(od...)
	return obj, diags
}

func stringListValue(ctx context.Context, ss []string) (types.List, diag.Diagnostics) {
	if ss == nil {
		ss = []string{}
	}
	return types.ListValueFrom(ctx, types.StringType, ss)
}

func familyToString(f int) string {
	switch f {
	case unix.AF_INET:
		return "inet"
	case unix.AF_INET6:
		return "inet6"
	case unix.AF_MPLS:
		return "mpls"
	default:
		return ""
	}
}

func encapTypeToString(t int) string {
	switch t {
	case unix.LWTUNNEL_ENCAP_MPLS:
		return "mpls"
	case unix.LWTUNNEL_ENCAP_IP:
		return "ip"
	case unix.LWTUNNEL_ENCAP_IP6:
		return "ip6"
	case unix.LWTUNNEL_ENCAP_ILA:
		return "ila"
	case unix.LWTUNNEL_ENCAP_BPF:
		return "bpf"
	case unix.LWTUNNEL_ENCAP_SEG6:
		return "seg6"
	case unix.LWTUNNEL_ENCAP_SEG6_LOCAL:
		return "seg6local"
	default:
		return fmt.Sprintf("%d", t)
	}
}
