package provider

import (
	"context"
	"encoding/hex"
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
	"golang.org/x/sys/unix"
)

var _ resource.Resource = &XfrmStateResource{}

type XfrmStateResource struct {
	client *netlinkClient.Client
}

func NewXfrmStateResource() resource.Resource { return &XfrmStateResource{} }

func (r *XfrmStateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_xfrm_state"
}

func (r *XfrmStateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an XFRM/IPsec security association state (ip xfrm state).",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"src":           schema.StringAttribute{Required: true, Description: "Source IP."},
			"dst":           schema.StringAttribute{Required: true, Description: "Destination IP."},
			"proto":         schema.StringAttribute{Required: true, Description: "Protocol (esp, ah, comp)."},
			"spi":           schema.Int64Attribute{Required: true, Description: "Security Parameter Index."},
			"mode":          schema.StringAttribute{Optional: true, Description: "Mode (transport, tunnel)."},
			"reqid":         schema.Int64Attribute{Optional: true, Description: "Request ID."},
			"replay_window": schema.Int64Attribute{Optional: true, Description: "Replay window size."},
			"mark":          schema.Int64Attribute{Optional: true, Description: "Mark value."},
			"mark_mask":     schema.Int64Attribute{Optional: true, Description: "Mark mask."},
			"if_id":         schema.Int64Attribute{Optional: true, Description: "Interface ID."},
			"family":        schema.StringAttribute{Computed: true, Description: "Address family."},
		},
		Blocks: map[string]schema.Block{
			"auth": schema.SingleNestedBlock{
				Description: "Authentication algorithm.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{Optional: true},
					"key":  schema.StringAttribute{Optional: true, Sensitive: true},
				},
			},
			"crypt": schema.SingleNestedBlock{
				Description: "Encryption algorithm.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{Optional: true},
					"key":  schema.StringAttribute{Optional: true, Sensitive: true},
				},
			},
			"aead": schema.SingleNestedBlock{
				Description: "AEAD algorithm.",
				Attributes: map[string]schema.Attribute{
					"name":    schema.StringAttribute{Optional: true},
					"key":     schema.StringAttribute{Optional: true, Sensitive: true},
					"icv_len": schema.Int64Attribute{Optional: true},
				},
			},
		},
	}
}

func (r *XfrmStateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *XfrmStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.XfrmStateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.buildState(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build xfrm state", err.Error())
		return
	}

	if err := r.client.XfrmStateAdd(state); err != nil {
		resp.Diagnostics.AddError("Failed to add xfrm state", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%d", data.Src.ValueString(), data.Dst.ValueString(), data.SPI.ValueInt64()))
	data.Family = types.StringValue(xfrmFamily(data.Src.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *XfrmStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.XfrmStateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := &vnl.XfrmState{
		Src:   net.ParseIP(data.Src.ValueString()),
		Dst:   net.ParseIP(data.Dst.ValueString()),
		Proto: xfrmProto(data.Proto.ValueString()),
		Spi:   int(data.SPI.ValueInt64()),
	}

	_, err := r.client.XfrmStateGet(state)
	if err != nil {
		resp.State.RemoveResource(ctx)
	}
}

func (r *XfrmStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.XfrmStateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.buildState(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build xfrm state", err.Error())
		return
	}

	if err := r.client.XfrmStateUpdate(state); err != nil {
		resp.Diagnostics.AddError("Failed to update xfrm state", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s|%s|%d", data.Src.ValueString(), data.Dst.ValueString(), data.SPI.ValueInt64()))
	data.Family = types.StringValue(xfrmFamily(data.Src.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *XfrmStateResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.XfrmStateModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := &vnl.XfrmState{
		Src:   net.ParseIP(data.Src.ValueString()),
		Dst:   net.ParseIP(data.Dst.ValueString()),
		Proto: xfrmProto(data.Proto.ValueString()),
		Spi:   int(data.SPI.ValueInt64()),
	}

	if err := r.client.XfrmStateDel(state); err != nil {
		resp.Diagnostics.AddError("Failed to delete xfrm state", err.Error())
	}
}

func (r *XfrmStateResource) buildState(data *models.XfrmStateModel) (*vnl.XfrmState, error) {
	state := &vnl.XfrmState{
		Src:   net.ParseIP(data.Src.ValueString()),
		Dst:   net.ParseIP(data.Dst.ValueString()),
		Proto: xfrmProto(data.Proto.ValueString()),
		Spi:   int(data.SPI.ValueInt64()),
		Mode:  vnl.XFRM_MODE_TUNNEL,
	}

	if !data.Mode.IsNull() && data.Mode.ValueString() == "transport" {
		state.Mode = vnl.XFRM_MODE_TRANSPORT
	}

	if !data.Reqid.IsNull() {
		state.Reqid = int(data.Reqid.ValueInt64())
	}

	if !data.ReplayWindow.IsNull() {
		state.ReplayWindow = int(data.ReplayWindow.ValueInt64())
	}

	if !data.Mark.IsNull() {
		state.Mark = &vnl.XfrmMark{
			Value: uint32(data.Mark.ValueInt64()),
		}
		if !data.MarkMask.IsNull() {
			state.Mark.Mask = uint32(data.MarkMask.ValueInt64())
		}
	}

	if !data.IfID.IsNull() {
		state.Ifid = int(data.IfID.ValueInt64())
	}

	if data.Auth != nil {
		key, err := hex.DecodeString(data.Auth.Key.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid auth key: %w", err)
		}
		state.Auth = &vnl.XfrmStateAlgo{
			Name: data.Auth.Name.ValueString(),
			Key:  key,
		}
	}

	if data.Crypt != nil {
		key, err := hex.DecodeString(data.Crypt.Key.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid crypt key: %w", err)
		}
		state.Crypt = &vnl.XfrmStateAlgo{
			Name: data.Crypt.Name.ValueString(),
			Key:  key,
		}
	}

	if data.Aead != nil {
		key, err := hex.DecodeString(data.Aead.Key.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid aead key: %w", err)
		}
		state.Aead = &vnl.XfrmStateAlgo{
			Name:   data.Aead.Name.ValueString(),
			Key:    key,
			ICVLen: int(data.Aead.ICVLen.ValueInt64()),
		}
	}

	return state, nil
}

func xfrmProto(s string) vnl.Proto {
	switch s {
	case "esp":
		return vnl.XFRM_PROTO_ESP
	case "ah":
		return vnl.XFRM_PROTO_AH
	case "comp":
		return vnl.XFRM_PROTO_COMP
	default:
		return vnl.XFRM_PROTO_ESP
	}
}

func xfrmFamily(ip string) string {
	if net.ParseIP(ip).To4() != nil {
		return "inet"
	}
	return "inet6"
}

// Ensure unix is used
var _ = unix.AF_INET
