package provider

import (
	"context"
	"fmt"
	"net"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &LinkDataSource{}

type LinkDataSource struct{ client *netlinkClient.Client }

func NewLinkDataSource() datasource.DataSource { return &LinkDataSource{} }

func (d *LinkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (d *LinkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read information about a network link.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"name":         schema.StringAttribute{Required: true, Description: "Interface name."},
			"type":         schema.StringAttribute{Computed: true},
			"mtu":          schema.Int64Attribute{Computed: true},
			"mac_address":  schema.StringAttribute{Computed: true},
			"admin_status": schema.StringAttribute{Computed: true},
			"oper_status":  schema.StringAttribute{Computed: true},
			"if_index":     schema.Int64Attribute{Computed: true},
			"tx_queue_len": schema.Int64Attribute{Computed: true},
			"master":       schema.StringAttribute{Computed: true},
		},
	}
}

func (d *LinkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netlinkClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected type", fmt.Sprintf("Expected *netlink.Client, got %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *LinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type linkDS struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Type        types.String `tfsdk:"type"`
		MTU         types.Int64  `tfsdk:"mtu"`
		MacAddress  types.String `tfsdk:"mac_address"`
		AdminStatus types.String `tfsdk:"admin_status"`
		OperStatus  types.String `tfsdk:"oper_status"`
		IfIndex     types.Int64  `tfsdk:"if_index"`
		TxQueueLen  types.Int64  `tfsdk:"tx_queue_len"`
		Master      types.String `tfsdk:"master"`
	}
	var data linkDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := d.client.LinkByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Link not found", err.Error())
		return
	}

	attrs := link.Attrs()
	data.ID = types.StringValue(attrs.Name)
	data.MTU = types.Int64Value(int64(attrs.MTU))
	data.IfIndex = types.Int64Value(int64(attrs.Index))
	data.TxQueueLen = types.Int64Value(int64(attrs.TxQLen))
	data.Type = types.StringValue(link.Type())

	if attrs.HardwareAddr != nil {
		data.MacAddress = types.StringValue(attrs.HardwareAddr.String())
	} else {
		data.MacAddress = types.StringValue("")
	}

	if attrs.Flags&net.FlagUp != 0 {
		data.AdminStatus = types.StringValue("up")
	} else {
		data.AdminStatus = types.StringValue("down")
	}
	data.OperStatus = types.StringValue(attrs.OperState.String())

	if attrs.MasterIndex > 0 {
		master, err := d.client.LinkByIndex(attrs.MasterIndex)
		if err == nil {
			data.Master = types.StringValue(master.Attrs().Name)
		} else {
			data.Master = types.StringNull()
		}
	} else {
		data.Master = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
