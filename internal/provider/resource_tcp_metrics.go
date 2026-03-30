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

var _ resource.Resource = &TcpMetricsResource{}

type TcpMetricsResource struct{ client *netlinkClient.Client }

func NewTcpMetricsResource() resource.Resource { return &TcpMetricsResource{} }

func (r *TcpMetricsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tcp_metrics"
}

func (r *TcpMetricsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages TCP metrics cache entries (ip tcp_metrics).",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"address":  schema.StringAttribute{Required: true, Description: "Destination address."},
			"rtt":      schema.Int64Attribute{Computed: true, Description: "RTT (microseconds)."},
			"rttvar":   schema.Int64Attribute{Computed: true, Description: "RTT variance."},
			"ssthresh": schema.Int64Attribute{Computed: true, Description: "Slow start threshold."},
			"cwnd":     schema.Int64Attribute{Computed: true, Description: "Congestion window."},
		},
	}
}

func (r *TcpMetricsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TcpMetricsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TcpMetricsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TCP metrics are typically auto-created by the kernel; this manages the cache entry
	data.ID = types.StringValue(data.Address.ValueString())
	data.RTT = types.Int64Value(0)
	data.RTTVar = types.Int64Value(0)
	data.SSTHRESH = types.Int64Value(0)
	data.CWND = types.Int64Value(0)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TcpMetricsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TcpMetricsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	output, err := r.client.TcpMetricsList()
	if err != nil {
		return
	}
	if !containsString(output, data.Address.ValueString()) {
		resp.State.RemoveResource(ctx)
	}
}

func (r *TcpMetricsResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "TCP metrics entries are managed by the kernel.")
}

func (r *TcpMetricsResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TcpMetricsModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Ignore errors - the entry may not exist in the kernel if no TCP connection was made
	_ = r.client.TcpMetricsDel(data.Address.ValueString())
}
