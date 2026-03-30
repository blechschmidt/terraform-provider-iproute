package provider

import (
	"context"
	"fmt"

	netlinkPkg "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &NetnsDataSource{}

type NetnsDataSource struct{}

func NewNetnsDataSource() datasource.DataSource { return &NetnsDataSource{} }

func (d *NetnsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_netns"
}

func (d *NetnsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List network namespaces.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"namespaces": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of namespace names."},
		},
	}
}

func (d *NetnsDataSource) Configure(_ context.Context, _ datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {}

func (d *NetnsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	type netnsDS struct {
		ID         types.String `tfsdk:"id"`
		Namespaces types.List   `tfsdk:"namespaces"`
	}

	names, err := netlinkPkg.ListNamespaces()
	if err != nil {
		resp.Diagnostics.AddError("Failed to list namespaces", err.Error())
		return
	}

	nsList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)

	data := netnsDS{
		ID:         types.StringValue("namespaces"),
		Namespaces: nsList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Ensure fmt is used
var _ = fmt.Sprintf
