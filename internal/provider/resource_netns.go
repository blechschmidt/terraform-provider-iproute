package provider

import (
	"context"
	"fmt"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &NetnsResource{}
	_ resource.ResourceWithImportState = &NetnsResource{}
)

type NetnsResource struct{}

func NewNetnsResource() resource.Resource {
	return &NetnsResource{}
}

func (r *NetnsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_netns"
}

func (r *NetnsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a network namespace (ip netns).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Namespace name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *NetnsResource) Configure(_ context.Context, _ resource.ConfigureRequest, _ *resource.ConfigureResponse) {}

func (r *NetnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.NetnsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	if err := netlinkPkg.CreateNamespace(name); err != nil {
		resp.Diagnostics.AddError("Failed to create namespace", err.Error())
		return
	}

	data.ID = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.NetnsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exists, err := netlinkPkg.NamespaceExists(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to check namespace", err.Error())
		return
	}
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetnsResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported",
		"Network namespaces cannot be updated, only replaced.")
}

func (r *NetnsResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.NetnsModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := netlinkPkg.DeleteNamespace(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete namespace", err.Error())
	}
}

func (r *NetnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// Ensure fmt is used
var _ = fmt.Sprintf
