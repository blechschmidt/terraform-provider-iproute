package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
)

var _ resource.Resource = &XfrmPolicyResource{}

type XfrmPolicyResource struct {
	client *netlinkClient.Client
}

func NewXfrmPolicyResource() resource.Resource { return &XfrmPolicyResource{} }

func (r *XfrmPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_xfrm_policy"
}

func (r *XfrmPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an XFRM/IPsec policy (ip xfrm policy).",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"src":       schema.StringAttribute{Required: true, Description: "Source prefix."},
			"dst":       schema.StringAttribute{Required: true, Description: "Destination prefix."},
			"dir":       schema.StringAttribute{Required: true, Description: "Direction (in, out, fwd)."},
			"priority":  schema.Int64Attribute{Optional: true, Description: "Policy priority."},
			"action":    schema.StringAttribute{Optional: true, Description: "Action (allow, block)."},
			"proto":     schema.StringAttribute{Optional: true, Description: "Protocol."},
			"src_port":  schema.Int64Attribute{Optional: true, Description: "Source port."},
			"dst_port":  schema.Int64Attribute{Optional: true, Description: "Destination port."},
			"mark":      schema.Int64Attribute{Optional: true, Description: "Mark value."},
			"mark_mask": schema.Int64Attribute{Optional: true, Description: "Mark mask."},
			"if_id":     schema.Int64Attribute{Optional: true, Description: "Interface ID."},
			"family":    schema.StringAttribute{Computed: true, Description: "Address family."},
			"templates": schema.ListNestedAttribute{
				Optional: true, Description: "Policy templates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"src":   schema.StringAttribute{Optional: true},
						"dst":   schema.StringAttribute{Optional: true},
						"proto": schema.StringAttribute{Optional: true},
						"mode":  schema.StringAttribute{Optional: true},
						"reqid": schema.Int64Attribute{Optional: true},
						"spi":   schema.Int64Attribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *XfrmPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *XfrmPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.XfrmPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.buildPolicy(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build xfrm policy", err.Error())
		return
	}

	if err := r.client.XfrmPolicyAdd(policy); err != nil {
		resp.Diagnostics.AddError("Failed to add xfrm policy", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s", data.Src.ValueString(), data.Dst.ValueString(), data.Dir.ValueString()))
	data.Family = types.StringValue(xfrmFamily(data.Src.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *XfrmPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.XfrmPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, _ := r.buildPolicy(ctx, &data)
	if policy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	_, err := r.client.XfrmPolicyGet(policy)
	if err != nil {
		resp.State.RemoveResource(ctx)
	}
}

func (r *XfrmPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.XfrmPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.buildPolicy(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build xfrm policy", err.Error())
		return
	}

	if err := r.client.XfrmPolicyUpdate(policy); err != nil {
		resp.Diagnostics.AddError("Failed to update xfrm policy", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%s", data.Src.ValueString(), data.Dst.ValueString(), data.Dir.ValueString()))
	data.Family = types.StringValue(xfrmFamily(data.Src.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *XfrmPolicyResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.XfrmPolicyModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, _ := r.buildPolicy(context.Background(), &data)
	if policy == nil {
		return
	}

	if err := r.client.XfrmPolicyDel(policy); err != nil {
		resp.Diagnostics.AddError("Failed to delete xfrm policy", err.Error())
	}
}

func (r *XfrmPolicyResource) buildPolicy(ctx context.Context, data *models.XfrmPolicyModel) (*vnl.XfrmPolicy, error) {
	_, srcNet, err := net.ParseCIDR(data.Src.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid src: %w", err)
	}
	_, dstNet, err := net.ParseCIDR(data.Dst.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid dst: %w", err)
	}

	policy := &vnl.XfrmPolicy{
		Src: srcNet,
		Dst: dstNet,
		Dir: xfrmDir(data.Dir.ValueString()),
	}

	if !data.Priority.IsNull() {
		policy.Priority = int(data.Priority.ValueInt64())
	}

	if !data.Action.IsNull() && data.Action.ValueString() == "block" {
		policy.Action = vnl.XFRM_POLICY_BLOCK
	}

	if !data.Mark.IsNull() {
		policy.Mark = &vnl.XfrmMark{Value: uint32(data.Mark.ValueInt64())}
		if !data.MarkMask.IsNull() {
			policy.Mark.Mask = uint32(data.MarkMask.ValueInt64())
		}
	}

	if !data.IfID.IsNull() {
		policy.Ifid = int(data.IfID.ValueInt64())
	}

	// Templates
	if !data.Templates.IsNull() && !data.Templates.IsUnknown() {
		var templates []models.XfrmPolicyTemplate
		data.Templates.ElementsAs(ctx, &templates, false)
		for _, t := range templates {
			tmpl := vnl.XfrmPolicyTmpl{
				Proto: xfrmProto(t.Proto.ValueString()),
				Mode:  vnl.XFRM_MODE_TUNNEL,
			}
			if !t.Mode.IsNull() && t.Mode.ValueString() == "transport" {
				tmpl.Mode = vnl.XFRM_MODE_TRANSPORT
			}
			if !t.Src.IsNull() {
				tmpl.Src = net.ParseIP(t.Src.ValueString())
			}
			if !t.Dst.IsNull() {
				tmpl.Dst = net.ParseIP(t.Dst.ValueString())
			}
			if !t.Reqid.IsNull() {
				tmpl.Reqid = int(t.Reqid.ValueInt64())
			}
			if !t.SPI.IsNull() {
				tmpl.Spi = int(t.SPI.ValueInt64())
			}
			policy.Tmpls = append(policy.Tmpls, tmpl)
		}
	}

	return policy, nil
}

func xfrmDir(s string) vnl.Dir {
	switch s {
	case "in":
		return vnl.XFRM_DIR_IN
	case "out":
		return vnl.XFRM_DIR_OUT
	case "fwd":
		return vnl.XFRM_DIR_FWD
	default:
		return vnl.XFRM_DIR_OUT
	}
}
