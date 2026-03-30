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
	vnl "github.com/vishvananda/netlink"
)

var _ resource.Resource = &TuntapResource{}

type TuntapResource struct{ client *netlinkClient.Client }

func NewTuntapResource() resource.Resource { return &TuntapResource{} }

func (r *TuntapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tuntap"
}

func (r *TuntapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a TUN/TAP device (ip tuntap).",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":        schema.StringAttribute{Required: true, Description: "Device name."},
			"mode":        schema.StringAttribute{Required: true, Description: "Mode: tun or tap."},
			"owner":       schema.Int64Attribute{Optional: true, Description: "Owner UID."},
			"group":       schema.Int64Attribute{Optional: true, Description: "Group GID."},
			"multi_queue": schema.BoolAttribute{Optional: true, Description: "Enable multi-queue."},
			"persist":     schema.BoolAttribute{Optional: true, Description: "Persistent device."},
		},
	}
}

func (r *TuntapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TuntapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TuntapModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mode := vnl.TUNTAP_MODE_TUN
	if data.Mode.ValueString() == "tap" {
		mode = vnl.TUNTAP_MODE_TAP
	}

	tt := &vnl.Tuntap{
		LinkAttrs: vnl.LinkAttrs{Name: data.Name.ValueString()},
		Mode:      mode,
	}
	if !data.Owner.IsNull() { tt.Owner = uint32(data.Owner.ValueInt64()) }
	if !data.Group.IsNull() { tt.Group = uint32(data.Group.ValueInt64()) }
	if !data.MultiQueue.IsNull() && data.MultiQueue.ValueBool() {
		tt.Flags = vnl.TUNTAP_MULTI_QUEUE_DEFAULTS
	}

	if err := r.client.LinkAdd(tt); err != nil {
		resp.Diagnostics.AddError("Failed to create TUN/TAP device", err.Error())
		return
	}

	data.ID = types.StringValue(data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TuntapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TuntapModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
	}
}

func (r *TuntapResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "TUN/TAP devices must be replaced.")
}

func (r *TuntapResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TuntapModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	link, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		return
	}
	if err := r.client.LinkDel(link); err != nil {
		resp.Diagnostics.AddError("Failed to delete TUN/TAP device", err.Error())
	}
}
