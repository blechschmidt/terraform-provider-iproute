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

var _ datasource.DataSource = &NeighborDataSource{}

type NeighborDataSource struct{ client *netlinkClient.Client }

func NewNeighborDataSource() datasource.DataSource { return &NeighborDataSource{} }

func (d *NeighborDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_neighbor"
}

func (d *NeighborDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read neighbor/ARP entries.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"device":    schema.StringAttribute{Optional: true, Description: "Filter by interface."},
			"neighbors": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of neighbor entries."},
		},
	}
}

func (d *NeighborDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NeighborDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type neighDS struct {
		ID        types.String `tfsdk:"id"`
		Device    types.String `tfsdk:"device"`
		Neighbors types.List   `tfsdk:"neighbors"`
	}
	var data neighDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	linkIndex := 0
	if !data.Device.IsNull() {
		link, err := d.client.LinkByName(data.Device.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Interface not found", err.Error())
			return
		}
		linkIndex = link.Attrs().Index
	}

	neighs, err := d.client.NeighList(linkIndex, unix.AF_UNSPEC)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list neighbors", err.Error())
		return
	}

	entries := make([]string, 0, len(neighs))
	for _, n := range neighs {
		entry := n.IP.String()
		if n.HardwareAddr != nil {
			entry += " lladdr " + n.HardwareAddr.String()
		}
		entries = append(entries, entry)
	}

	neighList, diags := types.ListValueFrom(ctx, types.StringType, entries)
	resp.Diagnostics.Append(diags...)
	data.Neighbors = neighList
	data.ID = types.StringValue("neighbors")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
