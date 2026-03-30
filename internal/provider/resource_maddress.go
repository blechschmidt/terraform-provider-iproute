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

var _ resource.Resource = &MaddressResource{}

type MaddressResource struct{ client *netlinkClient.Client }

func NewMaddressResource() resource.Resource { return &MaddressResource{} }

func (r *MaddressResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maddress"
}

func (r *MaddressResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a multicast address (ip maddress).",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"device":  schema.StringAttribute{Required: true, Description: "Interface name."},
			"address": schema.StringAttribute{Required: true, Description: "Multicast address."},
		},
	}
}

func (r *MaddressResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.MaddressModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.MaddressAdd(data.Device.ValueString(), data.Address.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to add multicast address", err.Error())
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.Device.ValueString(), data.Address.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MaddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.MaddressModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	output, err := r.client.MaddressList(data.Device.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if !containsString(output, data.Address.ValueString()) {
		resp.State.RemoveResource(ctx)
	}
}

func (r *MaddressResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Multicast addresses must be replaced.")
}

func (r *MaddressResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.MaddressModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.MaddressDel(data.Device.ValueString(), data.Address.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete multicast address", err.Error())
	}
}
