package provider

import (
	"context"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &IprouteProvider{}

// IprouteProvider defines the provider implementation.
type IprouteProvider struct {
	version string
}

// IprouteProviderModel describes the provider data model.
type IprouteProviderModel struct {
	Namespace types.String `tfsdk:"namespace"`
}

// New returns a function that creates new provider instances.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IprouteProvider{
			version: version,
		}
	}
}

func (p *IprouteProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iproute"
	resp.Version = p.version
}

func (p *IprouteProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The iproute provider manages Linux networking resources using netlink. " +
			"It supports all resource types from iproute2 including links, addresses, routes, " +
			"rules, neighbors, and more.",
		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				Description: "Network namespace to operate in. If not set, operates in the default namespace.",
				Optional:    true,
			},
		},
	}
}

func (p *IprouteProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config IprouteProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namespace := ""
	if !config.Namespace.IsNull() && !config.Namespace.IsUnknown() {
		namespace = config.Namespace.ValueString()
	}

	client, err := netlinkClient.NewClient(namespace)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Netlink Client",
			"An unexpected error occurred when creating the netlink client. "+
				"Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *IprouteProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLinkResource,
		NewAddressResource,
		NewRouteResource,
		NewRuleResource,
		NewNeighborResource,
		NewNetnsResource,
		NewNexthopResource,
		NewTunnelResource,
		NewL2tpTunnelResource,
		NewL2tpSessionResource,
		NewFouResource,
		NewXfrmStateResource,
		NewXfrmPolicyResource,
		NewMacsecResource,
		NewTuntapResource,
		NewMaddressResource,
		NewTokenResource,
		NewTcpMetricsResource,
		NewSrResource,
	}
}

func (p *IprouteProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLinkDataSource,
		NewAddressDataSource,
		NewRouteDataSource,
		NewNeighborDataSource,
		NewRuleDataSource,
		NewNetnsDataSource,
		NewNexthopDataSource,
	}
}
