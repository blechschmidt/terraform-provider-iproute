package provider

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/example/terraform-provider-iproute/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	_ resource.Resource                = &AddressResource{}
	_ resource.ResourceWithImportState = &AddressResource{}
)

type AddressResource struct {
	client *netlinkClient.Client
}

func NewAddressResource() resource.Resource {
	return &AddressResource{}
}

func (r *AddressResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address"
}

func (r *AddressResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IP address on a network interface (ip address). RFC 8344.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "IP address in CIDR notation (e.g., 10.0.0.1/24).",
				Validators:  []validator.String{validators.IsCIDR()},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"device": schema.StringAttribute{
				Required:    true,
				Description: "Interface name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"peer": schema.StringAttribute{
				Optional:    true,
				Description: "Peer address for point-to-point interfaces.",
			},
			"broadcast": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Broadcast address.",
			},
			"label": schema.StringAttribute{
				Optional:    true,
				Description: "Address label.",
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Address scope (global, link, host, site).",
			},
			"family": schema.StringAttribute{
				Computed:    true,
				Description: "Address family (inet or inet6).",
			},
			"origin": schema.StringAttribute{
				Computed:    true,
				Description: "Address origin. RFC 8344.",
			},
			"flags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Address flags.",
			},
			"preferred_lifetime": schema.Int64Attribute{
				Optional:    true,
				Description: "Preferred lifetime in seconds.",
			},
			"valid_lifetime": schema.Int64Attribute{
				Optional:    true,
				Description: "Valid lifetime in seconds.",
			},
		},
	}
}

func (r *AddressResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.AddressModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Interface not found", err.Error())
		return
	}

	addr, err := vnl.ParseAddr(data.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid address", err.Error())
		return
	}

	if !data.Peer.IsNull() && !data.Peer.IsUnknown() {
		peer, err := vnl.ParseAddr(data.Peer.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid peer address", err.Error())
			return
		}
		addr.Peer = peer.IPNet
	}

	if !data.Broadcast.IsNull() && !data.Broadcast.IsUnknown() {
		addr.Broadcast = net.ParseIP(data.Broadcast.ValueString())
	}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		addr.Label = data.Label.ValueString()
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		addr.Scope = parseScope(data.Scope.ValueString())
	}

	if !data.PreferredLifetime.IsNull() {
		addr.PreferedLft = int(data.PreferredLifetime.ValueInt64())
	}
	if !data.ValidLifetime.IsNull() {
		addr.ValidLft = int(data.ValidLifetime.ValueInt64())
	}

	if err := r.client.AddrAdd(link, addr); err != nil {
		resp.Diagnostics.AddError("Failed to add address", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.Device.ValueString(), data.Address.ValueString()))
	r.readAddress(ctx, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.AddressModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	targetAddr, err := vnl.ParseAddr(data.Address.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	family := unix.AF_INET
	if targetAddr.IP.To4() == nil {
		family = unix.AF_INET6
	}

	addrs, err := r.client.AddrList(link, family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list addresses", err.Error())
		return
	}

	found := false
	for _, a := range addrs {
		if a.IPNet.String() == targetAddr.IPNet.String() {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	r.readAddress(ctx, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.AddressModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Interface not found", err.Error())
		return
	}

	addr, err := vnl.ParseAddr(data.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid address", err.Error())
		return
	}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		addr.Label = data.Label.ValueString()
	}

	if !data.PreferredLifetime.IsNull() {
		addr.PreferedLft = int(data.PreferredLifetime.ValueInt64())
	}
	if !data.ValidLifetime.IsNull() {
		addr.ValidLft = int(data.ValidLifetime.ValueInt64())
	}

	if err := r.client.AddrReplace(link, addr); err != nil {
		resp.Diagnostics.AddError("Failed to update address", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.Device.ValueString(), data.Address.ValueString()))
	r.readAddress(ctx, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddressResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.AddressModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Device.ValueString())
	if err != nil {
		return // interface already gone
	}

	addr, err := vnl.ParseAddr(data.Address.ValueString())
	if err != nil {
		return
	}

	if !data.Peer.IsNull() && !data.Peer.IsUnknown() {
		peer, err := vnl.ParseAddr(data.Peer.ValueString())
		if err == nil {
			addr.Peer = peer.IPNet
		}
	}

	if err := r.client.AddrDel(link, addr); err != nil {
		resp.Diagnostics.AddError("Failed to delete address", err.Error())
	}
}

func (r *AddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: device|address (e.g., eth0|10.0.0.1/24)
	parts := strings.SplitN(req.ID, "|", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Import ID must be in format: device|address (e.g., eth0|10.0.0.1/24)")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("address"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *AddressResource) readAddress(_ context.Context, data *models.AddressModel) {
	addr, _ := vnl.ParseAddr(data.Address.ValueString())
	if addr == nil {
		return
	}

	if addr.IP.To4() != nil {
		data.Family = types.StringValue("inet")
	} else {
		data.Family = types.StringValue("inet6")
	}

	data.Origin = types.StringValue("static")

	if data.Scope.IsNull() || data.Scope.IsUnknown() {
		data.Scope = types.StringValue("global")
	}

	if data.Flags.IsNull() || data.Flags.IsUnknown() {
		data.Flags = types.ListNull(types.StringType)
	}

	if data.Broadcast.IsNull() || data.Broadcast.IsUnknown() {
		data.Broadcast = types.StringNull()
	}
}

func parseScope(s string) int {
	switch strings.ToLower(s) {
	case "global":
		return unix.RT_SCOPE_UNIVERSE
	case "site":
		return unix.RT_SCOPE_SITE
	case "link":
		return unix.RT_SCOPE_LINK
	case "host":
		return unix.RT_SCOPE_HOST
	case "nowhere":
		return unix.RT_SCOPE_NOWHERE
	default:
		return unix.RT_SCOPE_UNIVERSE
	}
}
