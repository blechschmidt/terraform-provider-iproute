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

var _ resource.Resource = &MacsecResource{}

type MacsecResource struct{ client *netlinkClient.Client }

func NewMacsecResource() resource.Resource { return &MacsecResource{} }

func (r *MacsecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_macsec"
}

func (r *MacsecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a MACsec device (ip macsec).",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"parent":          schema.StringAttribute{Required: true, Description: "Parent interface."},
			"name":            schema.StringAttribute{Required: true, Description: "MACsec interface name."},
			"sci":             schema.StringAttribute{Optional: true, Description: "SCI value."},
			"port":            schema.Int64Attribute{Optional: true, Description: "Port number."},
			"encrypt":         schema.BoolAttribute{Optional: true, Description: "Enable encryption."},
			"cipher_suite":    schema.StringAttribute{Optional: true, Description: "Cipher suite."},
			"icv_len":         schema.Int64Attribute{Optional: true, Description: "ICV length."},
			"encoding_sa":     schema.Int64Attribute{Optional: true, Description: "Encoding SA."},
			"validate":        schema.StringAttribute{Optional: true, Description: "Validate mode."},
			"protect":         schema.BoolAttribute{Optional: true, Description: "Enable protection."},
			"replay_protect":  schema.BoolAttribute{Optional: true, Description: "Enable replay protection."},
			"window":          schema.Int64Attribute{Optional: true, Description: "Replay window size."},
		},
	}
}

func (r *MacsecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MacsecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.MacsecModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := &netlinkClient.Macsec{
		Parent: data.Parent.ValueString(),
		Name:   data.Name.ValueString(),
	}
	if !data.SCI.IsNull() { m.SCI = data.SCI.ValueString() }
	if !data.Port.IsNull() { m.Port = int(data.Port.ValueInt64()) }
	if !data.Encrypt.IsNull() { m.Encrypt = data.Encrypt.ValueBool() }
	if !data.CipherSuite.IsNull() { m.CipherSuite = data.CipherSuite.ValueString() }
	if !data.ICVLen.IsNull() { m.ICVLen = int(data.ICVLen.ValueInt64()) }
	if !data.Validate.IsNull() { m.Validate = data.Validate.ValueString() }
	if !data.Protect.IsNull() { m.Protect = data.Protect.ValueBool() }
	if !data.ReplayProtect.IsNull() {
		m.ReplayProtect = data.ReplayProtect.ValueBool()
		if !data.Window.IsNull() { m.Window = int(data.Window.ValueInt64()) }
	}

	if err := r.client.MacsecAdd(m); err != nil {
		resp.Diagnostics.AddError("Failed to create MACsec device", err.Error())
		return
	}

	data.ID = types.StringValue(data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MacsecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.MacsecModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
	}
}

func (r *MacsecResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "MACsec devices must be replaced.")
}

func (r *MacsecResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.MacsecModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.MacsecDel(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete MACsec device", err.Error())
	}
}
