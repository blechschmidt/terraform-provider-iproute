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

var _ resource.Resource = &TokenResource{}

type TokenResource struct{ client *netlinkClient.Client }

func NewTokenResource() resource.Resource { return &TokenResource{} }

func (r *TokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token"
}

func (r *TokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages IPv6 tokenized interface identifiers (ip token).",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"device": schema.StringAttribute{Required: true, Description: "Interface name."},
			"token":  schema.StringAttribute{Required: true, Description: "IPv6 interface token."},
		},
	}
}

func (r *TokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.TokenSet(data.Device.ValueString(), data.Token.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to set token", err.Error())
		return
	}
	data.ID = types.StringValue(data.Device.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	token, err := r.client.TokenGet(data.Device.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	data.Token = types.StringValue(token)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.TokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.TokenSet(data.Device.ValueString(), data.Token.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to set token", err.Error())
		return
	}
	data.ID = types.StringValue(data.Device.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TokenResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TokenModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Reset token to :: (all zeros)
	_ = r.client.TokenSet(data.Device.ValueString(), "::")
}
