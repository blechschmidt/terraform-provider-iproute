package provider

import (
	"context"
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &L2tpTunnelResource{}

type L2tpTunnelResource struct {
	client *netlinkClient.Client
}

func NewL2tpTunnelResource() resource.Resource { return &L2tpTunnelResource{} }

func (r *L2tpTunnelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2tp_tunnel"
}

func (r *L2tpTunnelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an L2TP tunnel (ip l2tp add tunnel).",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"tunnel_id":       schema.Int64Attribute{Required: true, Description: "Tunnel ID."},
			"peer_tunnel_id":  schema.Int64Attribute{Required: true, Description: "Peer tunnel ID."},
			"encap_type":      schema.StringAttribute{Required: true, Description: "Encapsulation type (udp or ip)."},
			"local":           schema.StringAttribute{Required: true, Description: "Local IP address."},
			"remote":          schema.StringAttribute{Required: true, Description: "Remote IP address."},
			"local_port":      schema.Int64Attribute{Optional: true, Description: "Local UDP port."},
			"remote_port":     schema.Int64Attribute{Optional: true, Description: "Remote UDP port."},
		},
	}
}

func (r *L2tpTunnelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *L2tpTunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.L2tpTunnelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t := &netlinkClient.L2tpTunnel{
		TunnelID:     int(data.TunnelID.ValueInt64()),
		PeerTunnelID: int(data.PeerTunnelID.ValueInt64()),
		EncapType:    data.EncapType.ValueString(),
		Local:        data.Local.ValueString(),
		Remote:       data.Remote.ValueString(),
	}
	if !data.LocalPort.IsNull() {
		t.LocalPort = int(data.LocalPort.ValueInt64())
	}
	if !data.RemotePort.IsNull() {
		t.RemotePort = int(data.RemotePort.ValueInt64())
	}

	if err := r.client.L2tpAddTunnel(t); err != nil {
		resp.Diagnostics.AddError("Failed to create L2TP tunnel", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d", data.TunnelID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2tpTunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.L2tpTunnelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// L2TP tunnels are checked by listing
	output, err := r.client.L2tpListTunnels()
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	tidStr := fmt.Sprintf("Tunnel %d", data.TunnelID.ValueInt64())
	if !containsString(output, tidStr) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2tpTunnelResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "L2TP tunnels must be replaced.")
}

func (r *L2tpTunnelResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.L2tpTunnelModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.L2tpDelTunnel(int(data.TunnelID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Failed to delete L2TP tunnel", err.Error())
	}
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && searchString(s, substr))
}

func searchString(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Ensure path import used
var _ = path.Root
