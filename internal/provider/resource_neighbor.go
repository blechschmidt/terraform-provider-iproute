package provider

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
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
	_ resource.Resource                = &NeighborResource{}
	_ resource.ResourceWithImportState = &NeighborResource{}
)

type NeighborResource struct {
	client *netlinkClient.Client
}

func NewNeighborResource() resource.Resource {
	return &NeighborResource{}
}

func (r *NeighborResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neighbor"
}

func (r *NeighborResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a neighbor/ARP entry (ip neigh). RFC 8344.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "IP address.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"lladdr": schema.StringAttribute{
				Required:    true,
				Description: "Link-layer address (MAC).",
			},
			"device": schema.StringAttribute{
				Required:    true,
				Description: "Interface name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Neighbor state (permanent, noarp, reachable, stale).",
			},
			"family": schema.StringAttribute{
				Computed:    true,
				Description: "Address family.",
			},
			"flags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Neighbor flags.",
			},
			"is_router": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether neighbor is a router. RFC 8344.",
			},
			"origin": schema.StringAttribute{
				Computed:    true,
				Description: "Neighbor origin. RFC 8344.",
			},
			"proxy": schema.BoolAttribute{
				Optional:    true,
				Description: "Add proxy ARP entry.",
			},
			"vni": schema.Int64Attribute{
				Optional:    true,
				Description: "VNI for FDB entries.",
			},
		},
	}
}

func (r *NeighborResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeighborResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.NeighborModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	neigh, err := r.buildNeigh(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build neighbor", err.Error())
		return
	}

	if err := r.client.NeighSet(neigh); err != nil {
		resp.Diagnostics.AddError("Failed to add neighbor", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.Device.ValueString(), data.Address.ValueString()))
	r.setComputed(&data, neigh)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeighborResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.NeighborModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	ip := net.ParseIP(data.Address.ValueString())
	family := unix.AF_INET
	if ip.To4() == nil {
		family = unix.AF_INET6
	}

	neighs, err := r.client.NeighList(link.Attrs().Index, family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list neighbors", err.Error())
		return
	}

	found := false
	for _, n := range neighs {
		if n.IP.Equal(ip) {
			data.LLAddr = types.StringValue(n.HardwareAddr.String())
			data.State = types.StringValue(neighStateToString(n.State))
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeighborResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.NeighborModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	neigh, err := r.buildNeigh(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build neighbor", err.Error())
		return
	}

	if err := r.client.NeighSet(neigh); err != nil {
		resp.Diagnostics.AddError("Failed to update neighbor", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.Device.ValueString(), data.Address.ValueString()))
	r.setComputed(&data, neigh)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeighborResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.NeighborModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	neigh, err := r.buildNeigh(&data)
	if err != nil {
		return
	}

	if err := r.client.NeighDel(neigh); err != nil {
		resp.Diagnostics.AddError("Failed to delete neighbor", err.Error())
	}
}

func (r *NeighborResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "|", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Import ID must be in format: device|address")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("address"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *NeighborResource) buildNeigh(data *models.NeighborModel) (*vnl.Neigh, error) {
	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		return nil, fmt.Errorf("device %q not found: %w", data.Device.ValueString(), err)
	}

	ip := net.ParseIP(data.Address.ValueString())
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", data.Address.ValueString())
	}

	hwAddr, err := net.ParseMAC(data.LLAddr.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address: %w", err)
	}

	family := unix.AF_INET
	if ip.To4() == nil {
		family = unix.AF_INET6
	}

	state := vnl.NUD_PERMANENT
	if !data.State.IsNull() && !data.State.IsUnknown() {
		state = parseNeighState(data.State.ValueString())
	}

	neigh := &vnl.Neigh{
		LinkIndex:    link.Attrs().Index,
		Family:       family,
		State:        state,
		IP:           ip,
		HardwareAddr: hwAddr,
	}

	if !data.Proxy.IsNull() && data.Proxy.ValueBool() {
		neigh.Flags = vnl.NTF_PROXY
	}

	return neigh, nil
}

func (r *NeighborResource) setComputed(data *models.NeighborModel, neigh *vnl.Neigh) {
	if neigh.Family == unix.AF_INET6 {
		data.Family = types.StringValue("inet6")
	} else {
		data.Family = types.StringValue("inet")
	}
	data.Origin = types.StringValue("static")
	data.IsRouter = types.BoolValue(false)
	if data.State.IsNull() || data.State.IsUnknown() {
		data.State = types.StringValue("permanent")
	}
	if data.Flags.IsNull() || data.Flags.IsUnknown() {
		data.Flags = types.ListNull(types.StringType)
	}
}

func parseNeighState(s string) int {
	switch strings.ToLower(s) {
	case "permanent":
		return vnl.NUD_PERMANENT
	case "noarp":
		return vnl.NUD_NOARP
	case "reachable":
		return vnl.NUD_REACHABLE
	case "stale":
		return vnl.NUD_STALE
	default:
		return vnl.NUD_PERMANENT
	}
}

func neighStateToString(s int) string {
	switch s {
	case vnl.NUD_PERMANENT:
		return "permanent"
	case vnl.NUD_NOARP:
		return "noarp"
	case vnl.NUD_REACHABLE:
		return "reachable"
	case vnl.NUD_STALE:
		return "stale"
	case vnl.NUD_INCOMPLETE:
		return "incomplete"
	case vnl.NUD_DELAY:
		return "delay"
	case vnl.NUD_PROBE:
		return "probe"
	case vnl.NUD_FAILED:
		return "failed"
	default:
		return fmt.Sprintf("%d", s)
	}
}
