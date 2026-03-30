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

var _ datasource.DataSource = &RouteDataSource{}

type RouteDataSource struct{ client *netlinkClient.Client }

func NewRouteDataSource() datasource.DataSource { return &RouteDataSource{} }

func (d *RouteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (d *RouteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read routing table entries.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"family": schema.StringAttribute{Optional: true, Description: "Address family (inet, inet6)."},
			"table":  schema.Int64Attribute{Optional: true, Description: "Routing table ID."},
			"routes": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of routes."},
		},
	}
}

func (d *RouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type routeDS struct {
		ID     types.String `tfsdk:"id"`
		Family types.String `tfsdk:"family"`
		Table  types.Int64  `tfsdk:"table"`
		Routes types.List   `tfsdk:"routes"`
	}
	var data routeDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	family := unix.AF_INET
	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		family = unix.AF_INET6
	}

	routes, err := d.client.RouteList(nil, family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list routes", err.Error())
		return
	}

	routeStrs := make([]string, 0, len(routes))
	for _, r := range routes {
		dst := "default"
		if r.Dst != nil {
			dst = r.Dst.String()
		}
		gw := ""
		if r.Gw != nil {
			gw = " via " + r.Gw.String()
		}
		routeStrs = append(routeStrs, dst+gw)
	}

	routeList, diags := types.ListValueFrom(ctx, types.StringType, routeStrs)
	resp.Diagnostics.Append(diags...)
	data.Routes = routeList
	data.ID = types.StringValue("routes")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
