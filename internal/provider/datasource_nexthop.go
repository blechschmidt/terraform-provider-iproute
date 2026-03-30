package provider

import (
	"context"
	"fmt"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &NexthopDataSource{}

type NexthopDataSource struct{ client *netlinkClient.Client }

func NewNexthopDataSource() datasource.DataSource { return &NexthopDataSource{} }

func (d *NexthopDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nexthop"
}

func (d *NexthopDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List nexthop objects.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true},
			"nexthops": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of nexthop entries."},
		},
	}
}

func (d *NexthopDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NexthopDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	type nhDS struct {
		ID       types.String `tfsdk:"id"`
		Nexthops types.List   `tfsdk:"nexthops"`
	}

	nhs, err := d.client.NexthopList()
	if err != nil {
		resp.Diagnostics.AddError("Failed to list nexthops", err.Error())
		return
	}

	var entries []string
	for _, nh := range nhs {
		entry := fmt.Sprintf("id %d", nh.ID)
		if nh.Gateway != nil {
			entry += " via " + nh.Gateway.String()
		}
		if nh.Blackhole {
			entry += " blackhole"
		}
		entries = append(entries, entry)
	}

	nhList, diags := types.ListValueFrom(ctx, types.StringType, entries)
	resp.Diagnostics.Append(diags...)

	data := nhDS{
		ID:       types.StringValue("nexthops"),
		Nexthops: nhList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
