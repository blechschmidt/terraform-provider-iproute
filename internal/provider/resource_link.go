package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/example/terraform-provider-iproute/internal/models"
	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/example/terraform-provider-iproute/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
)

var (
	_ resource.Resource                = &LinkResource{}
	_ resource.ResourceWithImportState = &LinkResource{}
)

type LinkResource struct {
	client *netlinkClient.Client
}

func NewLinkResource() resource.Resource {
	return &LinkResource{}
}

func (r *LinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (r *LinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a network link/interface (ip link). Supports 22 link types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Interface name (max 15 characters).",
				Validators:  []validator.String{validators.IsInterfaceName()},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Link type: dummy, bridge, veth, vlan, vxlan, bond, macvlan, macvtap, ipvlan, ifb, gre, gretap, ip6gre, ip6gretap, sit, vti, vti6, ip6tnl, ipip, geneve, wireguard, tuntap.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Interface description (maps to ifAlias, RFC 8343).",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Administrative state (up/down). Defaults to true.",
			},
			"mac_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Hardware MAC address.",
				Validators:  []validator.String{validators.IsMACAddress()},
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum transmission unit.",
			},
			"tx_queue_len": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Transmit queue length.",
			},
			"master": schema.StringAttribute{
				Optional:    true,
				Description: "Master device name (e.g., bridge name).",
			},
			"admin_status": schema.StringAttribute{
				Computed:    true,
				Description: "Administrative status (up/down). RFC 8343.",
			},
			"oper_status": schema.StringAttribute{
				Computed:    true,
				Description: "Operational status. RFC 8343.",
			},
			"if_index": schema.Int64Attribute{
				Computed:    true,
				Description: "Interface index.",
			},
			"speed": schema.Int64Attribute{
				Computed:    true,
				Description: "Interface speed in Mbps.",
			},
			"statistics": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Interface statistics (computed). RFC 8343.",
				Attributes: map[string]schema.Attribute{
					"rx_bytes":   schema.Int64Attribute{Computed: true},
					"tx_bytes":   schema.Int64Attribute{Computed: true},
					"rx_packets": schema.Int64Attribute{Computed: true},
					"tx_packets": schema.Int64Attribute{Computed: true},
					"rx_errors":  schema.Int64Attribute{Computed: true},
					"tx_errors":  schema.Int64Attribute{Computed: true},
					"rx_dropped": schema.Int64Attribute{Computed: true},
					"tx_dropped": schema.Int64Attribute{Computed: true},
					"multicast":  schema.Int64Attribute{Computed: true},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"bridge":    bridgeSchemaBlock(),
			"veth":      vethSchemaBlock(),
			"vlan":      vlanSchemaBlock(),
			"vxlan":     vxlanSchemaBlock(),
			"bond":      bondSchemaBlock(),
			"macvlan":   macvlanSchemaBlock(),
			"macvtap":   macvtapSchemaBlock(),
			"ipvlan":    ipvlanSchemaBlock(),
			"dummy":     dummySchemaBlock(),
			"ifb":       ifbSchemaBlock(),
			"gre":       greSchemaBlock(),
			"sit":       sitSchemaBlock(),
			"vti":       vtiSchemaBlock(),
			"ip6tnl":    ip6tnlSchemaBlock(),
			"ipip":      ipipSchemaBlock(),
			"geneve":    geneveSchemaBlock(),
			"wireguard": wireguardSchemaBlock(),
			"tuntap":    tuntapSchemaBlock(),
		},
	}
}

func (r *LinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.LinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.buildLink(&data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build link", err.Error())
		return
	}

	if err := r.client.LinkAdd(link); err != nil {
		resp.Diagnostics.AddError("Failed to create link", err.Error())
		return
	}

	// For veth, the peer is auto-created
	// Set properties that can't be set during creation
	created, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to find created link", err.Error())
		return
	}

	// Set MTU if specified
	if !data.MTU.IsNull() && !data.MTU.IsUnknown() {
		if err := r.client.LinkSetMTU(created, int(data.MTU.ValueInt64())); err != nil {
			resp.Diagnostics.AddError("Failed to set MTU", err.Error())
			return
		}
	}

	// Set TxQueueLen if specified
	if !data.TxQueueLen.IsNull() && !data.TxQueueLen.IsUnknown() {
		if err := r.client.LinkSetTxQLen(created, int(data.TxQueueLen.ValueInt64())); err != nil {
			resp.Diagnostics.AddError("Failed to set TxQueueLen", err.Error())
			return
		}
	}

	// Set MAC address if specified
	if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() {
		hwAddr, err := net.ParseMAC(data.MacAddress.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse MAC address", err.Error())
			return
		}
		if err := r.client.LinkSetHardwareAddr(created, hwAddr); err != nil {
			resp.Diagnostics.AddError("Failed to set MAC address", err.Error())
			return
		}
	}

	// Set alias/description
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		if err := r.client.LinkSetAlias(created, data.Description.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set description", err.Error())
			return
		}
	}

	// Set master
	if !data.Master.IsNull() && !data.Master.IsUnknown() {
		master, err := r.client.LinkByName(data.Master.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to find master device", err.Error())
			return
		}
		if err := r.client.LinkSetMaster(created, master); err != nil {
			resp.Diagnostics.AddError("Failed to set master", err.Error())
			return
		}
	}

	// Set up/down
	if data.Enabled.ValueBool() {
		if err := r.client.LinkSetUp(created); err != nil {
			resp.Diagnostics.AddError("Failed to bring link up", err.Error())
			return
		}
	}

	// Read back the link state
	r.readLinkInto(&data)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.LinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		// Link no longer exists
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateModel(&data, link)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.LinkModel
	var state models.LinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to find link", err.Error())
		return
	}

	// Update name if changed
	if data.Name.ValueString() != state.Name.ValueString() {
		if err := r.client.LinkSetName(link, data.Name.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to rename link", err.Error())
			return
		}
		link, err = r.client.LinkByName(data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to find renamed link", err.Error())
			return
		}
	}

	// Update MTU
	if !data.MTU.IsNull() && !data.MTU.IsUnknown() && data.MTU.ValueInt64() != state.MTU.ValueInt64() {
		if err := r.client.LinkSetMTU(link, int(data.MTU.ValueInt64())); err != nil {
			resp.Diagnostics.AddError("Failed to set MTU", err.Error())
			return
		}
	}

	// Update TxQueueLen
	if !data.TxQueueLen.IsNull() && !data.TxQueueLen.IsUnknown() && data.TxQueueLen.ValueInt64() != state.TxQueueLen.ValueInt64() {
		if err := r.client.LinkSetTxQLen(link, int(data.TxQueueLen.ValueInt64())); err != nil {
			resp.Diagnostics.AddError("Failed to set TxQueueLen", err.Error())
			return
		}
	}

	// Update MAC address
	if !data.MacAddress.IsNull() && !data.MacAddress.IsUnknown() && data.MacAddress.ValueString() != state.MacAddress.ValueString() {
		hwAddr, err := net.ParseMAC(data.MacAddress.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse MAC address", err.Error())
			return
		}
		if err := r.client.LinkSetHardwareAddr(link, hwAddr); err != nil {
			resp.Diagnostics.AddError("Failed to set MAC address", err.Error())
			return
		}
	}

	// Update description
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		if err := r.client.LinkSetAlias(link, data.Description.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set description", err.Error())
			return
		}
	}

	// Update master
	if !data.Master.Equal(state.Master) {
		if data.Master.IsNull() {
			if err := r.client.LinkSetNoMaster(link); err != nil {
				resp.Diagnostics.AddError("Failed to remove master", err.Error())
				return
			}
		} else {
			master, err := r.client.LinkByName(data.Master.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Failed to find master device", err.Error())
				return
			}
			if err := r.client.LinkSetMaster(link, master); err != nil {
				resp.Diagnostics.AddError("Failed to set master", err.Error())
				return
			}
		}
	}

	// Update enabled state
	if data.Enabled.ValueBool() != state.Enabled.ValueBool() {
		if data.Enabled.ValueBool() {
			if err := r.client.LinkSetUp(link); err != nil {
				resp.Diagnostics.AddError("Failed to bring link up", err.Error())
				return
			}
		} else {
			if err := r.client.LinkSetDown(link); err != nil {
				resp.Diagnostics.AddError("Failed to bring link down", err.Error())
				return
			}
		}
	}

	r.readLinkInto(&data)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.LinkModel
	resp.Diagnostics.Append(req.State.Get(context.Background(), &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		// Already gone
		return
	}

	if err := r.client.LinkDel(link); err != nil {
		resp.Diagnostics.AddError("Failed to delete link", err.Error())
	}
}

func (r *LinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	// Set name from ID for read
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *LinkResource) readLinkInto(data *models.LinkModel) {
	link, err := r.client.LinkByName(data.Name.ValueString())
	if err != nil {
		return
	}
	r.populateModel(data, link)
}

func (r *LinkResource) populateModel(data *models.LinkModel, link vnl.Link) {
	attrs := link.Attrs()
	data.ID = types.StringValue(attrs.Name)
	data.Name = types.StringValue(attrs.Name)
	data.IfIndex = types.Int64Value(int64(attrs.Index))
	data.MTU = types.Int64Value(int64(attrs.MTU))
	data.TxQueueLen = types.Int64Value(int64(attrs.TxQLen))

	if attrs.HardwareAddr != nil {
		data.MacAddress = types.StringValue(attrs.HardwareAddr.String())
	} else {
		data.MacAddress = types.StringValue("")
	}

	if attrs.Flags&net.FlagUp != 0 {
		data.Enabled = types.BoolValue(true)
		data.AdminStatus = types.StringValue("up")
	} else {
		data.Enabled = types.BoolValue(false)
		data.AdminStatus = types.StringValue("down")
	}

	data.OperStatus = types.StringValue(attrs.OperState.String())
	data.Speed = types.Int64Value(0)

	if attrs.Alias != "" {
		data.Description = types.StringValue(attrs.Alias)
	}

	if attrs.MasterIndex > 0 {
		master, err := r.client.LinkByIndex(attrs.MasterIndex)
		if err == nil {
			data.Master = types.StringValue(master.Attrs().Name)
		}
	}

	// Set type based on link type
	switch link.(type) {
	case *vnl.Bridge:
		data.Type = types.StringValue("bridge")
	case *vnl.Veth:
		data.Type = types.StringValue("veth")
	case *vnl.Vlan:
		data.Type = types.StringValue("vlan")
	case *vnl.Vxlan:
		data.Type = types.StringValue("vxlan")
	case *vnl.Bond:
		data.Type = types.StringValue("bond")
	case *vnl.Macvlan:
		data.Type = types.StringValue("macvlan")
	case *vnl.Macvtap:
		data.Type = types.StringValue("macvtap")
	case *vnl.IPVlan:
		data.Type = types.StringValue("ipvlan")
	case *vnl.Dummy:
		data.Type = types.StringValue("dummy")
	case *vnl.Ifb:
		data.Type = types.StringValue("ifb")
	case *vnl.Gretun:
		data.Type = types.StringValue("gre")
	case *vnl.Gretap:
		data.Type = types.StringValue("gretap")
	case *vnl.Iptun:
		data.Type = types.StringValue("ipip")
	case *vnl.Sittun:
		data.Type = types.StringValue("sit")
	case *vnl.Vti:
		data.Type = types.StringValue("vti")
	case *vnl.Ip6tnl:
		data.Type = types.StringValue("ip6tnl")
	case *vnl.Geneve:
		data.Type = types.StringValue("geneve")
	case *vnl.Wireguard:
		data.Type = types.StringValue("wireguard")
	case *vnl.Tuntap:
		data.Type = types.StringValue("tuntap")
	default:
		if data.Type.IsNull() || data.Type.IsUnknown() {
			data.Type = types.StringValue("unknown")
		}
	}

	// Statistics
	statsAttrTypes := map[string]attr.Type{
		"rx_bytes":   types.Int64Type,
		"tx_bytes":   types.Int64Type,
		"rx_packets": types.Int64Type,
		"tx_packets": types.Int64Type,
		"rx_errors":  types.Int64Type,
		"tx_errors":  types.Int64Type,
		"rx_dropped": types.Int64Type,
		"tx_dropped": types.Int64Type,
		"multicast":  types.Int64Type,
	}
	if attrs.Statistics != nil {
		data.Statistics, _ = types.ObjectValue(statsAttrTypes, map[string]attr.Value{
			"rx_bytes":   types.Int64Value(int64(attrs.Statistics.RxBytes)),
			"tx_bytes":   types.Int64Value(int64(attrs.Statistics.TxBytes)),
			"rx_packets": types.Int64Value(int64(attrs.Statistics.RxPackets)),
			"tx_packets": types.Int64Value(int64(attrs.Statistics.TxPackets)),
			"rx_errors":  types.Int64Value(int64(attrs.Statistics.RxErrors)),
			"tx_errors":  types.Int64Value(int64(attrs.Statistics.TxErrors)),
			"rx_dropped": types.Int64Value(int64(attrs.Statistics.RxDropped)),
			"tx_dropped": types.Int64Value(int64(attrs.Statistics.TxDropped)),
			"multicast":  types.Int64Value(int64(attrs.Statistics.Multicast)),
		})
	} else {
		data.Statistics = types.ObjectNull(statsAttrTypes)
	}
}

func (r *LinkResource) buildLink(data *models.LinkModel) (vnl.Link, error) {
	name := data.Name.ValueString()
	linkType := data.Type.ValueString()

	switch linkType {
	case "dummy":
		return &vnl.Dummy{LinkAttrs: vnl.LinkAttrs{Name: name}}, nil
	case "bridge":
		return buildBridge(name, data.Bridge)
	case "veth":
		return buildVeth(name, data.Veth)
	case "vlan":
		return buildVlan(name, data.Vlan, r.client)
	case "vxlan":
		return buildVxlan(name, data.Vxlan)
	case "bond":
		return buildBond(name, data.Bond)
	case "macvlan":
		return buildMacvlan(name, data.Macvlan, r.client)
	case "macvtap":
		return buildMacvtap(name, data.Macvtap, r.client)
	case "ipvlan":
		return buildIpvlan(name, data.Ipvlan, r.client)
	case "ifb":
		return &vnl.Ifb{LinkAttrs: vnl.LinkAttrs{Name: name}}, nil
	case "gre":
		return buildGre(name, data.Gre)
	case "gretap":
		return buildGretap(name, data.Gre)
	case "sit":
		return buildSit(name, data.Sit)
	case "vti":
		return buildVti(name, data.Vti)
	case "ip6tnl":
		return buildIp6tnl(name, data.Ip6tnl)
	case "ipip":
		return buildIpip(name, data.Ipip)
	case "geneve":
		return buildGeneve(name, data.Geneve)
	case "wireguard":
		return &vnl.Wireguard{LinkAttrs: vnl.LinkAttrs{Name: name}}, nil
	case "tuntap":
		return buildTuntap(name, data.Tuntap)
	default:
		return nil, fmt.Errorf("unsupported link type: %s", linkType)
	}
}
