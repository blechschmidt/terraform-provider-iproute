package provider

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	_ resource.Resource                = &RouteResource{}
	_ resource.ResourceWithImportState = &RouteResource{}
)

type RouteResource struct {
	client *netlinkClient.Client
}

func NewRouteResource() resource.Resource {
	return &RouteResource{}
}

func (r *RouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *RouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a routing table entry (ip route).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "Destination in CIDR notation or 'default'.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Description: "Gateway IP address.",
			},
			"device": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Output device name.",
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "Source address for route.",
			},
			"metric": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Route metric/priority.",
			},
			"table": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Routing table ID.",
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Route scope.",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Route protocol.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Route type (unicast, blackhole, unreachable, prohibit).",
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Description: "Route MTU.",
			},
			"advmss": schema.Int64Attribute{
				Optional:    true,
				Description: "Advertised MSS.",
			},
			"family": schema.StringAttribute{
				Computed:    true,
				Description: "Address family.",
			},
			"nexthop_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Nexthop ID.",
			},
			"onlink": schema.BoolAttribute{
				Optional:    true,
				Description: "Gateway is on-link even if not in any connected subnet.",
			},
			"multipath": schema.ListNestedAttribute{
				Optional:    true,
				Description: "ECMP multipath routes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gateway": schema.StringAttribute{Optional: true},
						"device":  schema.StringAttribute{Optional: true},
						"weight":  schema.Int64Attribute{Optional: true},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"encap": schema.SingleNestedBlock{
				Description: "Route encapsulation.",
				Attributes: map[string]schema.Attribute{
					"type":   schema.StringAttribute{Optional: true, Description: "Encap type (mpls, seg6, etc)."},
					"labels": schema.ListAttribute{ElementType: types.StringType, Optional: true, Description: "Labels."},
				},
			},
		},
	}
}

func (r *RouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netlinkClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netlink.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.RouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.buildRoute(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build route", err.Error())
		return
	}

	if err := r.client.RouteAdd(route); err != nil {
		resp.Diagnostics.AddError("Failed to add route", err.Error())
		return
	}

	r.setID(&data)
	r.setComputed(&data, route)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.RouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.buildRoute(&data)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	family := unix.AF_INET
	if route.Dst != nil && route.Dst.IP.To4() == nil {
		family = unix.AF_INET6
	}

	filterMask := uint64(vnl.RT_FILTER_DST)
	if route.Table > 0 {
		filterMask |= uint64(vnl.RT_FILTER_TABLE)
	}

	routes, err := r.client.RouteListFiltered(family, route, filterMask)
	if err != nil || len(routes) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	r.setComputed(&data, &routes[0])
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.RouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.buildRoute(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build route", err.Error())
		return
	}

	if err := r.client.RouteReplace(route); err != nil {
		resp.Diagnostics.AddError("Failed to replace route", err.Error())
		return
	}

	r.setID(&data)
	r.setComputed(&data, route)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouteResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.RouteModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.buildRoute(&data)
	if err != nil {
		return
	}

	if err := r.client.RouteDel(route); err != nil {
		resp.Diagnostics.AddError("Failed to delete route", err.Error())
	}
}

func (r *RouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: "destination|table" (e.g., "10.1.0.0/24|main")
	parts := strings.SplitN(req.ID, "|", 2)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination"), parts[0])...)
	if len(parts) > 1 && parts[1] != "main" {
		table, err := strconv.ParseInt(parts[1], 10, 64)
		if err == nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("table"), table)...)
		}
	}
}

func (r *RouteResource) buildRoute(data *models.RouteModel) (*vnl.Route, error) {
	route := &vnl.Route{}

	dest := data.Destination.ValueString()
	if dest == "default" {
		route.Dst = nil
	} else {
		_, dst, err := net.ParseCIDR(dest)
		if err != nil {
			return nil, fmt.Errorf("invalid destination %q: %w", dest, err)
		}
		route.Dst = dst
	}

	if !data.Gateway.IsNull() && !data.Gateway.IsUnknown() {
		gw := net.ParseIP(data.Gateway.ValueString())
		if gw == nil {
			return nil, fmt.Errorf("invalid gateway %q", data.Gateway.ValueString())
		}
		route.Gw = gw
	}

	if !data.Device.IsNull() && !data.Device.IsUnknown() && data.Device.ValueString() != "" {
		link, err := r.client.LinkByName(data.Device.ValueString())
		if err != nil {
			return nil, fmt.Errorf("device %q not found: %w", data.Device.ValueString(), err)
		}
		route.LinkIndex = link.Attrs().Index
	}

	if !data.Source.IsNull() && !data.Source.IsUnknown() {
		route.Src = net.ParseIP(data.Source.ValueString())
	}

	if !data.Metric.IsNull() && !data.Metric.IsUnknown() {
		route.Priority = int(data.Metric.ValueInt64())
	}

	if !data.Table.IsNull() && !data.Table.IsUnknown() {
		route.Table = int(data.Table.ValueInt64())
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		route.Scope = vnl.Scope(parseScope(data.Scope.ValueString()))
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		switch data.Type.ValueString() {
		case "blackhole":
			route.Type = unix.RTN_BLACKHOLE
		case "unreachable":
			route.Type = unix.RTN_UNREACHABLE
		case "prohibit":
			route.Type = unix.RTN_PROHIBIT
		default:
			route.Type = unix.RTN_UNICAST
		}
	}

	if !data.Onlink.IsNull() && data.Onlink.ValueBool() {
		route.Flags |= int(vnl.FLAG_ONLINK)
	}

	return route, nil
}

func (r *RouteResource) setID(data *models.RouteModel) {
	dest := data.Destination.ValueString()
	table := "main"
	if !data.Table.IsNull() && !data.Table.IsUnknown() {
		table = fmt.Sprintf("%d", data.Table.ValueInt64())
	}
	data.ID = types.StringValue(fmt.Sprintf("%s|%s", dest, table))
}

func (r *RouteResource) setComputed(data *models.RouteModel, route *vnl.Route) {
	if route.Dst != nil {
		if route.Dst.IP.To4() != nil {
			data.Family = types.StringValue("inet")
		} else {
			data.Family = types.StringValue("inet6")
		}
	} else {
		data.Family = types.StringValue("inet")
	}

	if data.Table.IsNull() || data.Table.IsUnknown() {
		data.Table = types.Int64Value(int64(route.Table))
	}
	if data.Metric.IsNull() || data.Metric.IsUnknown() {
		data.Metric = types.Int64Value(int64(route.Priority))
	}
	if data.Scope.IsNull() || data.Scope.IsUnknown() {
		data.Scope = types.StringValue(scopeToString(int(route.Scope)))
	}
	if data.Protocol.IsNull() || data.Protocol.IsUnknown() {
		data.Protocol = types.StringValue(protocolToString(route.Protocol))
	}
	if data.Type.IsNull() || data.Type.IsUnknown() {
		data.Type = types.StringValue(routeTypeToString(route.Type))
	}

	if data.Device.IsNull() || data.Device.IsUnknown() {
		if route.LinkIndex > 0 {
			link, err := r.client.LinkByIndex(route.LinkIndex)
			if err == nil {
				data.Device = types.StringValue(link.Attrs().Name)
			}
		} else {
			data.Device = types.StringValue("")
		}
	}

	if data.Multipath.IsNull() || data.Multipath.IsUnknown() {
		data.Multipath = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"device":  types.StringType,
				"gateway": types.StringType,
				"weight":  types.Int64Type,
			},
		})
	}
}

func scopeToString(s int) string {
	switch s {
	case unix.RT_SCOPE_UNIVERSE:
		return "global"
	case unix.RT_SCOPE_SITE:
		return "site"
	case unix.RT_SCOPE_LINK:
		return "link"
	case unix.RT_SCOPE_HOST:
		return "host"
	case unix.RT_SCOPE_NOWHERE:
		return "nowhere"
	default:
		return fmt.Sprintf("%d", s)
	}
}

func protocolToString(p vnl.RouteProtocol) string {
	switch p {
	case unix.RTPROT_BOOT:
		return "boot"
	case unix.RTPROT_KERNEL:
		return "kernel"
	case unix.RTPROT_STATIC:
		return "static"
	default:
		return fmt.Sprintf("%d", p)
	}
}

func routeTypeToString(t int) string {
	switch t {
	case unix.RTN_UNICAST:
		return "unicast"
	case unix.RTN_BLACKHOLE:
		return "blackhole"
	case unix.RTN_UNREACHABLE:
		return "unreachable"
	case unix.RTN_PROHIBIT:
		return "prohibit"
	default:
		return "unicast"
	}
}
