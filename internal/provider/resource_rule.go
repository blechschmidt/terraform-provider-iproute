package provider

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	_ resource.Resource                = &RuleResource{}
	_ resource.ResourceWithImportState = &RuleResource{}
)

type RuleResource struct {
	client *netlinkClient.Client
}

func NewRuleResource() resource.Resource {
	return &RuleResource{}
}

func (r *RuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule"
}

func (r *RuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a routing policy rule (ip rule).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"priority":           schema.Int64Attribute{Optional: true, Computed: true, Description: "Rule priority."},
			"family":             schema.StringAttribute{Optional: true, Computed: true, Description: "Address family."},
			"src":                schema.StringAttribute{Optional: true, Description: "Source prefix."},
			"dst":                schema.StringAttribute{Optional: true, Description: "Destination prefix."},
			"iif_name":           schema.StringAttribute{Optional: true, Description: "Input interface."},
			"oif_name":           schema.StringAttribute{Optional: true, Description: "Output interface."},
			"table":              schema.Int64Attribute{Optional: true, Computed: true, Description: "Routing table."},
			"fwmark":             schema.Int64Attribute{Optional: true, Description: "Firewall mark."},
			"fwmask":             schema.Int64Attribute{Optional: true, Description: "Firewall mark mask."},
			"tos":                schema.Int64Attribute{Optional: true, Description: "TOS value."},
			"action":             schema.StringAttribute{Optional: true, Computed: true, Description: "Rule action."},
			"goto_priority":      schema.Int64Attribute{Optional: true, Description: "Goto priority."},
			"suppress_prefix_len": schema.Int64Attribute{Optional: true, Description: "Suppress prefix length."},
			"suppress_if_group":  schema.Int64Attribute{Optional: true, Description: "Suppress if group."},
			"ip_proto":           schema.Int64Attribute{Optional: true, Description: "IP protocol."},
			"sport_range":        schema.StringAttribute{Optional: true, Description: "Source port range (start-end)."},
			"dport_range":        schema.StringAttribute{Optional: true, Description: "Destination port range (start-end)."},
			"uid_range":          schema.StringAttribute{Optional: true, Description: "UID range (start-end)."},
			"invert":             schema.BoolAttribute{Optional: true, Description: "Invert match."},
		},
	}
}

func (r *RuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netlinkClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netlink.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *RuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.RuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := r.buildRule(&data)

	if err := r.client.RuleAdd(rule); err != nil {
		resp.Diagnostics.AddError("Failed to add rule", err.Error())
		return
	}

	r.setID(&data)
	r.setComputed(&data, rule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.RuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	family := unix.AF_INET
	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		family = unix.AF_INET6
	}

	rules, err := r.client.RuleList(family)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list rules", err.Error())
		return
	}

	targetRule := r.buildRule(&data)
	var foundRule *vnl.Rule
	for _, rl := range rules {
		if r.rulesMatch(&rl, targetRule) {
			foundRule = &rl
			break
		}
	}

	if foundRule == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.setComputed(&data, foundRule)
	if foundRule.Src != nil && (data.Src.IsNull() || data.Src.IsUnknown()) {
		data.Src = types.StringValue(foundRule.Src.String())
	}
	if foundRule.Dst != nil && (data.Dst.IsNull() || data.Dst.IsUnknown()) {
		data.Dst = types.StringValue(foundRule.Dst.String())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.RuleModel
	var state models.RuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete old rule
	oldRule := r.buildRule(&state)
	_ = r.client.RuleDel(oldRule)

	// Add new rule
	newRule := r.buildRule(&data)
	if err := r.client.RuleAdd(newRule); err != nil {
		resp.Diagnostics.AddError("Failed to add rule", err.Error())
		return
	}

	r.setID(&data)
	r.setComputed(&data, newRule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.RuleModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := r.buildRule(&data)
	if err := r.client.RuleDel(rule); err != nil {
		resp.Diagnostics.AddError("Failed to delete rule", err.Error())
	}
}

func (r *RuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	// Parse priority from the import ID
	priority, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be the rule priority number")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("priority"), priority)...)
}

func (r *RuleResource) buildRule(data *models.RuleModel) *vnl.Rule {
	rule := vnl.NewRule()

	if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
		rule.Priority = int(data.Priority.ValueInt64())
	}

	if !data.Family.IsNull() && data.Family.ValueString() == "inet6" {
		rule.Family = unix.AF_INET6
	} else {
		rule.Family = unix.AF_INET
	}

	if !data.Src.IsNull() && !data.Src.IsUnknown() {
		_, src, err := net.ParseCIDR(data.Src.ValueString())
		if err == nil {
			rule.Src = src
		}
	}

	if !data.Dst.IsNull() && !data.Dst.IsUnknown() {
		_, dst, err := net.ParseCIDR(data.Dst.ValueString())
		if err == nil {
			rule.Dst = dst
		}
	}

	if !data.IifName.IsNull() && !data.IifName.IsUnknown() {
		rule.IifName = data.IifName.ValueString()
	}

	if !data.OifName.IsNull() && !data.OifName.IsUnknown() {
		rule.OifName = data.OifName.ValueString()
	}

	if !data.Table.IsNull() && !data.Table.IsUnknown() {
		rule.Table = int(data.Table.ValueInt64())
	}

	if !data.FwMark.IsNull() && !data.FwMark.IsUnknown() {
		rule.Mark = uint32(data.FwMark.ValueInt64())
	}

	if !data.FwMask.IsNull() && !data.FwMask.IsUnknown() {
		v := uint32(data.FwMask.ValueInt64())
		rule.Mask = &v
	}

	if !data.TosMatch.IsNull() && !data.TosMatch.IsUnknown() {
		rule.Tos = uint(data.TosMatch.ValueInt64())
	}

	if !data.Goto.IsNull() && !data.Goto.IsUnknown() {
		rule.Goto = int(data.Goto.ValueInt64())
	}

	if !data.SuppressPrefixLen.IsNull() && !data.SuppressPrefixLen.IsUnknown() {
		rule.SuppressPrefixlen = int(data.SuppressPrefixLen.ValueInt64())
	}

	if !data.SuppressIfGroup.IsNull() && !data.SuppressIfGroup.IsUnknown() {
		rule.SuppressIfgroup = int(data.SuppressIfGroup.ValueInt64())
	}

	if !data.Invert.IsNull() && data.Invert.ValueBool() {
		rule.Invert = true
	}

	return rule
}

func (r *RuleResource) rulesMatch(a, b *vnl.Rule) bool {
	if a.Priority != b.Priority {
		return false
	}
	if b.Table != 0 && a.Table != b.Table {
		return false
	}
	if a.Src != nil && b.Src != nil && a.Src.String() != b.Src.String() {
		return false
	}
	if a.Dst != nil && b.Dst != nil && a.Dst.String() != b.Dst.String() {
		return false
	}
	return true
}

func (r *RuleResource) setID(data *models.RuleModel) {
	priority := "0"
	if !data.Priority.IsNull() {
		priority = fmt.Sprintf("%d", data.Priority.ValueInt64())
	}
	data.ID = types.StringValue(priority)
}

func (r *RuleResource) setComputed(data *models.RuleModel, rule *vnl.Rule) {
	if data.Priority.IsNull() || data.Priority.IsUnknown() {
		data.Priority = types.Int64Value(int64(rule.Priority))
	}
	if data.Family.IsNull() || data.Family.IsUnknown() {
		if rule.Family == unix.AF_INET6 {
			data.Family = types.StringValue("inet6")
		} else {
			data.Family = types.StringValue("inet")
		}
	}
	if data.Table.IsNull() || data.Table.IsUnknown() {
		data.Table = types.Int64Value(int64(rule.Table))
	}
	if data.Action.IsNull() || data.Action.IsUnknown() {
		data.Action = types.StringValue("lookup")
	}
}
