package provider

import (
	"context"
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &L2tpSessionResource{}

type L2tpSessionResource struct {
	client *netlinkClient.Client
}

func NewL2tpSessionResource() resource.Resource { return &L2tpSessionResource{} }

func (r *L2tpSessionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2tp_session"
}

func (r *L2tpSessionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an L2TP session (ip l2tp add session).",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"tunnel_id":        schema.Int64Attribute{Required: true, Description: "Tunnel ID."},
			"session_id":       schema.Int64Attribute{Required: true, Description: "Session ID."},
			"peer_session_id":  schema.Int64Attribute{Required: true, Description: "Peer session ID."},
			"name":             schema.StringAttribute{Optional: true, Description: "Session interface name."},
			"cookie":           schema.StringAttribute{Optional: true, Description: "Cookie value."},
			"peer_cookie":      schema.StringAttribute{Optional: true, Description: "Peer cookie value."},
			"l2spec_type":      schema.StringAttribute{Optional: true, Description: "L2-specific header type."},
		},
	}
}

func (r *L2tpSessionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *L2tpSessionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.L2tpSessionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &netlinkClient.L2tpSession{
		TunnelID:      int(data.TunnelID.ValueInt64()),
		SessionID:     int(data.SessionID.ValueInt64()),
		PeerSessionID: int(data.PeerSessionID.ValueInt64()),
	}
	if !data.Name.IsNull() {
		s.Name = data.Name.ValueString()
	}
	if !data.Cookie.IsNull() {
		s.Cookie = data.Cookie.ValueString()
	}
	if !data.PeerCookie.IsNull() {
		s.PeerCookie = data.PeerCookie.ValueString()
	}
	if !data.L2specType.IsNull() {
		s.L2specType = data.L2specType.ValueString()
	}

	if err := r.client.L2tpAddSession(s); err != nil {
		resp.Diagnostics.AddError("Failed to create L2TP session", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d|%d", data.TunnelID.ValueInt64(), data.SessionID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *L2tpSessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.L2tpSessionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	output, err := r.client.L2tpListSessions()
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	sidStr := fmt.Sprintf("Session %d", data.SessionID.ValueInt64())
	if !containsString(output, sidStr) {
		resp.State.RemoveResource(ctx)
	}
}

func (r *L2tpSessionResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "L2TP sessions must be replaced.")
}

func (r *L2tpSessionResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.L2tpSessionModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.L2tpDelSession(int(data.TunnelID.ValueInt64()), int(data.SessionID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Failed to delete L2TP session", err.Error())
	}
}
