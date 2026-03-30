package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/sys/unix"
)

var (
	_ resource.Resource                = &NexthopResource{}
	_ resource.ResourceWithImportState = &NexthopResource{}
)

type NexthopResource struct {
	client *netlinkClient.Client
}

func NewNexthopResource() resource.Resource { return &NexthopResource{} }

func (r *NexthopResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nexthop"
}

func (r *NexthopResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a nexthop object (ip nexthop).",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"nhid":      schema.Int64Attribute{Required: true, Description: "Nexthop ID."},
			"gateway":   schema.StringAttribute{Optional: true, Description: "Gateway IP."},
			"device":    schema.StringAttribute{Optional: true, Description: "Output device."},
			"blackhole": schema.BoolAttribute{Optional: true, Description: "Blackhole nexthop."},
			"family":    schema.StringAttribute{Optional: true, Computed: true, Description: "Address family."},
			"group": schema.ListNestedAttribute{
				Optional: true, Description: "Nexthop group members.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.Int64Attribute{Required: true},
						"weight": schema.Int64Attribute{Optional: true},
					},
				},
			},
			"resilient": schema.BoolAttribute{Optional: true, Description: "Resilient nexthop group."},
			"fdb":       schema.BoolAttribute{Optional: true, Description: "FDB nexthop."},
		},
	}
}

func (r *NexthopResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NexthopResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.NexthopModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nh, err := r.buildNexthop(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build nexthop", err.Error())
		return
	}

	if err := r.client.NexthopAdd(nh); err != nil {
		resp.Diagnostics.AddError("Failed to add nexthop", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d", data.NhID.ValueInt64()))
	r.setComputed(&data, nh)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NexthopResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.NexthopModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nh, err := r.client.NexthopGet(int(data.NhID.ValueInt64()))
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.setComputed(&data, nh)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NexthopResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.NexthopModel
	var state models.NexthopModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete old, create new
	oldNh := &netlinkClient.Nexthop{ID: int(state.NhID.ValueInt64())}
	_ = r.client.NexthopDel(oldNh)

	nh, err := r.buildNexthop(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build nexthop", err.Error())
		return
	}

	if err := r.client.NexthopAdd(nh); err != nil {
		resp.Diagnostics.AddError("Failed to add nexthop", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d", data.NhID.ValueInt64()))
	r.setComputed(&data, nh)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NexthopResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.NexthopModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nh := &netlinkClient.Nexthop{ID: int(data.NhID.ValueInt64())}
	if err := r.client.NexthopDel(nh); err != nil {
		resp.Diagnostics.AddError("Failed to delete nexthop", err.Error())
	}
}

func (r *NexthopResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NexthopResource) buildNexthop(ctx context.Context, data *models.NexthopModel) (*netlinkClient.Nexthop, error) {
	nh := &netlinkClient.Nexthop{
		ID: int(data.NhID.ValueInt64()),
	}

	if !data.Gateway.IsNull() && !data.Gateway.IsUnknown() {
		nh.Gateway = net.ParseIP(data.Gateway.ValueString())
	}

	if !data.Device.IsNull() && !data.Device.IsUnknown() {
		link, err := r.client.LinkByName(data.Device.ValueString())
		if err != nil {
			return nil, fmt.Errorf("device %q not found: %w", data.Device.ValueString(), err)
		}
		nh.LinkIndex = link.Attrs().Index
	}

	if !data.Blackhole.IsNull() && data.Blackhole.ValueBool() {
		nh.Blackhole = true
	}

	if !data.FDB.IsNull() && data.FDB.ValueBool() {
		nh.FDB = true
	}

	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		nh.Family = unix.AF_INET6
	} else {
		nh.Family = unix.AF_INET
	}

	// Group members
	if !data.Group.IsNull() && !data.Group.IsUnknown() {
		var members []models.NexthopGroupMemberModel
		data.Group.ElementsAs(ctx, &members, false)
		for _, m := range members {
			nh.Group = append(nh.Group, netlinkClient.NexthopGroupMember{
				ID:     int(m.ID.ValueInt64()),
				Weight: int(m.Weight.ValueInt64()),
			})
		}
	}

	return nh, nil
}

func (r *NexthopResource) setComputed(data *models.NexthopModel, nh *netlinkClient.Nexthop) {
	if data.Family.IsNull() || data.Family.IsUnknown() {
		if nh.Family == unix.AF_INET6 {
			data.Family = types.StringValue("inet6")
		} else {
			data.Family = types.StringValue("inet")
		}
	}
	if data.Group.IsNull() || data.Group.IsUnknown() {
		data.Group = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":     types.Int64Type,
				"weight": types.Int64Type,
			},
		})
	}
}
