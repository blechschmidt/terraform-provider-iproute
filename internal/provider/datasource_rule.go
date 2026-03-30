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

var _ datasource.DataSource = &RuleDataSource{}

type RuleDataSource struct{ client *netlinkClient.Client }

func NewRuleDataSource() datasource.DataSource { return &RuleDataSource{} }

func (d *RuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule"
}

func (d *RuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read routing policy rules.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"family": schema.StringAttribute{Optional: true, Description: "Address family (inet, inet6)."},
			"rules":  schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "List of rules."},
		},
	}
}

func (d *RuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type ruleDS struct {
		ID     types.String `tfsdk:"id"`
		Family types.String `tfsdk:"family"`
		Rules  types.List   `tfsdk:"rules"`
	}
	var data ruleDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	family := unix.AF_INET
	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		family = unix.AF_INET6
	}

	rules, err := d.client.RuleList(family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list rules", err.Error())
		return
	}

	var ruleStrs []string
	for _, r := range rules {
		entry := fmt.Sprintf("%d:", r.Priority)
		if r.Src != nil {
			entry += " from " + r.Src.String()
		}
		if r.Dst != nil {
			entry += " to " + r.Dst.String()
		}
		entry += fmt.Sprintf(" lookup %d", r.Table)
		ruleStrs = append(ruleStrs, entry)
	}

	ruleList, diags := types.ListValueFrom(ctx, types.StringType, ruleStrs)
	resp.Diagnostics.Append(diags...)
	data.Rules = ruleList
	data.ID = types.StringValue("rules")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
