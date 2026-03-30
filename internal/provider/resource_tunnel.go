package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
)

var _ resource.Resource = &TunnelResource{}

type TunnelResource struct {
	client *netlinkClient.Client
}

func NewTunnelResource() resource.Resource { return &TunnelResource{} }

func (r *TunnelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunnel"
}

func (r *TunnelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IP tunnel (ip tunnel). Wraps link types gre/sit/ipip.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":   schema.StringAttribute{Required: true, Description: "Tunnel interface name."},
			"mode":   schema.StringAttribute{Required: true, Description: "Tunnel mode: gre, sit, ipip.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"local":  schema.StringAttribute{Optional: true, Description: "Local endpoint."},
			"remote": schema.StringAttribute{Optional: true, Description: "Remote endpoint."},
			"ttl":    schema.Int64Attribute{Optional: true, Description: "TTL."},
			"tos":    schema.Int64Attribute{Optional: true, Description: "TOS."},
		},
	}
}

func (r *TunnelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netlinkClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected type", fmt.Sprintf("Expected *netlink.Client, got %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *TunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var name, mode, local, remote types.String
	var ttl, tos types.Int64
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("mode"), &mode)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("local"), &local)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("remote"), &remote)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("ttl"), &ttl)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tos"), &tos)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.buildTunnel(name.ValueString(), mode.ValueString(), local, remote, ttl, tos)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build tunnel", err.Error())
		return
	}

	if err := r.client.LinkAdd(link); err != nil {
		resp.Diagnostics.AddError("Failed to create tunnel", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), name.ValueString())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mode"), mode)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("local"), local)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("remote"), remote)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ttl"), ttl)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tos"), tos)...)
}

func (r *TunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.LinkByName(name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
	}
}

func (r *TunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Tunnels must be replaced to change parameters.")
}

func (r *TunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(name.ValueString())
	if err != nil {
		return
	}
	if err := r.client.LinkDel(link); err != nil {
		resp.Diagnostics.AddError("Failed to delete tunnel", err.Error())
	}
}

func (r *TunnelResource) buildTunnel(name, mode string, local, remote types.String, ttl, tos types.Int64) (vnl.Link, error) {
	var localIP, remoteIP net.IP
	if !local.IsNull() {
		localIP = net.ParseIP(local.ValueString())
	}
	if !remote.IsNull() {
		remoteIP = net.ParseIP(remote.ValueString())
	}
	var ttlVal, tosVal uint8
	if !ttl.IsNull() {
		ttlVal = uint8(ttl.ValueInt64())
	}
	if !tos.IsNull() {
		tosVal = uint8(tos.ValueInt64())
	}

	switch mode {
	case "gre":
		return &vnl.Gretun{LinkAttrs: vnl.LinkAttrs{Name: name}, Local: localIP, Remote: remoteIP, Ttl: ttlVal, Tos: tosVal}, nil
	case "sit":
		return &vnl.Sittun{LinkAttrs: vnl.LinkAttrs{Name: name}, Local: localIP, Remote: remoteIP, Ttl: ttlVal, Tos: tosVal}, nil
	case "ipip":
		return &vnl.Iptun{LinkAttrs: vnl.LinkAttrs{Name: name}, Local: localIP, Remote: remoteIP, Ttl: ttlVal, Tos: tosVal}, nil
	default:
		return nil, fmt.Errorf("unsupported tunnel mode: %s", mode)
	}
}

// Ensure models import is used
var _ = models.RouteModel{}
