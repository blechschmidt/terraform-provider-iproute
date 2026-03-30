package provider

import (
	"context"
	"fmt"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/sys/unix"
)

var _ datasource.DataSource = &AddressDataSource{}

type AddressDataSource struct{ client *netlinkClient.Client }

func NewAddressDataSource() datasource.DataSource { return &AddressDataSource{} }

func (d *AddressDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address"
}

func (d *AddressDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read IP addresses on an interface.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"device":    schema.StringAttribute{Required: true, Description: "Interface name."},
			"family":    schema.StringAttribute{Optional: true, Description: "Address family filter (inet, inet6)."},
			"addresses": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of addresses in CIDR notation."},
		},
	}
}

func (d *AddressDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AddressDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type addrDS struct {
		ID        types.String `tfsdk:"id"`
		Device    types.String `tfsdk:"device"`
		Family    types.String `tfsdk:"family"`
		Addresses types.List   `tfsdk:"addresses"`
	}
	var data addrDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := d.client.LinkByName(data.Device.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Interface not found", err.Error())
		return
	}

	family := unix.AF_UNSPEC
	if !data.Family.IsNull() {
		if data.Family.ValueString() == "inet" {
			family = unix.AF_INET
		} else if data.Family.ValueString() == "inet6" {
			family = unix.AF_INET6
		}
	}

	addrs, err := d.client.AddrList(link, family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list addresses", err.Error())
		return
	}

	var addrStrs []string
	for _, a := range addrs {
		addrStrs = append(addrStrs, a.IPNet.String())
	}

	addrList, diags := types.ListValueFrom(ctx, types.StringType, addrStrs)
	resp.Diagnostics.Append(diags...)
	data.Addresses = addrList
	data.ID = types.StringValue(data.Device.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
