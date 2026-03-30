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
	"golang.org/x/sys/unix"
)

var _ resource.Resource = &FouResource{}

type FouResource struct {
	client *netlinkClient.Client
}

func NewFouResource() resource.Resource { return &FouResource{} }

func (r *FouResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fou"
}

func (r *FouResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Foo-over-UDP (FOU) receive port (ip fou).",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"port":        schema.Int64Attribute{Required: true, Description: "UDP port number."},
			"family":      schema.StringAttribute{Optional: true, Computed: true, Description: "Address family."},
			"protocol":    schema.Int64Attribute{Optional: true, Description: "IP protocol number."},
			"encap_type":  schema.StringAttribute{Optional: true, Description: "Encapsulation type: direct or gue."},
			"remote_port": schema.Int64Attribute{Optional: true, Description: "Remote port."},
			"local":       schema.StringAttribute{Optional: true, Description: "Local address."},
			"peer":        schema.StringAttribute{Optional: true, Description: "Peer address."},
		},
	}
}

func (r *FouResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FouResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.FouModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f := r.buildFou(&data)
	if err := r.client.FouAdd(f); err != nil {
		resp.Diagnostics.AddError("Failed to create FOU", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d", data.Port.ValueInt64()))
	if data.Family.IsNull() || data.Family.IsUnknown() {
		data.Family = types.StringValue("inet")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FouResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.FouModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	family := unix.AF_INET
	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		family = unix.AF_INET6
	}

	fous, err := r.client.FouList(family)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	port := int(data.Port.ValueInt64())
	found := false
	for _, f := range fous {
		if f.Port == port {
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
	}
}

func (r *FouResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "FOU entries must be replaced.")
}

func (r *FouResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.FouModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	f := r.buildFou(&data)
	if err := r.client.FouDel(f); err != nil {
		resp.Diagnostics.AddError("Failed to delete FOU", err.Error())
	}
}

func (r *FouResource) buildFou(data *models.FouModel) vnl.Fou {
	f := vnl.Fou{
		Port:   int(data.Port.ValueInt64()),
		Family: unix.AF_INET,
	}
	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		f.Family = unix.AF_INET6
	}
	if !data.Protocol.IsNull() {
		f.Protocol = int(data.Protocol.ValueInt64())
	}
	if !data.EncapType.IsNull() && data.EncapType.ValueString() == "gue" {
		f.EncapType = vnl.FOU_ENCAP_GUE
	} else {
		f.EncapType = vnl.FOU_ENCAP_DIRECT
	}
	return f
}
