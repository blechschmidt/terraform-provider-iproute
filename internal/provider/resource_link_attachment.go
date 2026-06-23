package provider

import (
	"context"
	"fmt"
	"net"

	netlinkClient "github.com/example/terraform-provider-iproute/internal/netlink"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vnl "github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

var (
	_ resource.Resource                = &LinkAttachmentResource{}
	_ resource.ResourceWithImportState = &LinkAttachmentResource{}
)

// LinkAttachmentResource manages the placement of a pre-existing network
// interface: it can move the interface into another network namespace, enslave
// it to a bridge and bring it up, without ever creating or destroying the
// interface itself. This is the declarative replacement for the imperative
// `ip link set <dev> netns <ns> master <br> up` sequence (e.g. bridging a
// physical host NIC, or a Docker-managed veth, into a bridge).
//
// On destroy it reverses the placement: the interface is returned to the
// provider's namespace, detached from its master and brought down, leaving it
// as the kernel would after the owning namespace went away.
type LinkAttachmentResource struct {
	client *netlinkClient.Client
}

func NewLinkAttachmentResource() resource.Resource {
	return &LinkAttachmentResource{}
}

type LinkAttachmentModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Netns  types.String `tfsdk:"netns"`
	Master types.String `tfsdk:"master"`
	Up     types.Bool   `tfsdk:"up"`
}

func (r *LinkAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link_attachment"
}

func (r *LinkAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attaches a pre-existing network interface to a namespace and/or bridge " +
			"without creating it. Moves the interface into the target network namespace " +
			"(optional), enslaves it to a master bridge (optional) and brings it up. " +
			"Reverses all of this on destroy. Use it to bridge a physical NIC or a " +
			"container veth into a bridge declaratively.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:      true,
				Description:   "Name of the existing interface to attach.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"netns": schema.StringAttribute{
				Optional: true,
				Description: "Target network namespace to move the interface into, using the " +
					"same syntax as the provider's namespace argument (name, name:, pid:, " +
					"path: or docker:). When unset, the interface stays in the provider's " +
					"namespace and is only enslaved/brought up there.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"master": schema.StringAttribute{
				Optional:    true,
				Description: "Master bridge to enslave the interface to (resolved in the target namespace).",
			},
			"up": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Administrative state. Defaults to true.",
			},
		},
	}
}

func (r *LinkAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// sourceNs returns the namespace handle the provider operates in, opening a
// handle to the current namespace when the provider uses the default one.
func (r *LinkAttachmentResource) sourceNs() (netns.NsHandle, bool, error) {
	if r.client.NsHandle >= 0 {
		return r.client.NsHandle, false, nil
	}
	h, err := netns.Get()
	return h, true, err
}

// targetHandle resolves the target namespace and returns a netlink handle in
// it plus the namespace handle. When netns is unset, it returns the provider's
// own handle and namespace.
func (r *LinkAttachmentResource) targetHandle(netnsSpec string) (*vnl.Handle, netns.NsHandle, bool, error) {
	if netnsSpec == "" {
		ns, owned, err := r.sourceNs()
		if err != nil {
			return nil, -1, false, err
		}
		return r.client.Handle, ns, owned, nil
	}
	ns, err := netlinkClient.ResolveNamespace(netnsSpec)
	if err != nil {
		return nil, -1, false, err
	}
	h, err := vnl.NewHandleAt(ns)
	if err != nil {
		ns.Close()
		return nil, -1, false, err
	}
	return h, ns, true, nil
}

func (r *LinkAttachmentResource) apply(data *LinkAttachmentModel) error {
	name := data.Name.ValueString()
	netnsSpec := data.Netns.ValueString()

	// Locate the interface in the provider's (source) namespace.
	link, err := r.client.Handle.LinkByName(name)
	movedAlready := false
	if err != nil {
		// It may already be in the target namespace (idempotent re-apply).
		if netnsSpec == "" {
			return fmt.Errorf("interface %q not found: %w", name, err)
		}
		movedAlready = true
	}

	tgtHandle, tgtNs, tgtOwned, err := r.targetHandle(netnsSpec)
	if err != nil {
		return fmt.Errorf("resolving target namespace: %w", err)
	}
	if tgtOwned {
		defer tgtNs.Close()
		if tgtHandle != r.client.Handle {
			defer tgtHandle.Close()
		}
	}

	// Move into the target namespace if needed.
	if netnsSpec != "" && !movedAlready {
		if err := r.client.Handle.LinkSetDown(link); err != nil {
			return fmt.Errorf("setting %q down before move: %w", name, err)
		}
		if err := r.client.Handle.LinkSetNsFd(link, int(tgtNs)); err != nil {
			return fmt.Errorf("moving %q into target namespace: %w", name, err)
		}
	}

	// Re-resolve in the target namespace.
	tgtLink, err := tgtHandle.LinkByName(name)
	if err != nil {
		return fmt.Errorf("interface %q not found in target namespace: %w", name, err)
	}

	// Enslave to master bridge.
	if !data.Master.IsNull() && data.Master.ValueString() != "" {
		master, err := tgtHandle.LinkByName(data.Master.ValueString())
		if err != nil {
			return fmt.Errorf("master %q not found in target namespace: %w", data.Master.ValueString(), err)
		}
		if err := tgtHandle.LinkSetMaster(tgtLink, master); err != nil {
			return fmt.Errorf("enslaving %q to %q: %w", name, data.Master.ValueString(), err)
		}
	} else {
		// No master desired: detach if currently enslaved.
		if tgtLink.Attrs().MasterIndex > 0 {
			if err := tgtHandle.LinkSetNoMaster(tgtLink); err != nil {
				return fmt.Errorf("detaching %q from master: %w", name, err)
			}
		}
	}

	// Administrative state.
	if data.Up.ValueBool() {
		if err := tgtHandle.LinkSetUp(tgtLink); err != nil {
			return fmt.Errorf("bringing %q up: %w", name, err)
		}
	} else {
		if err := tgtHandle.LinkSetDown(tgtLink); err != nil {
			return fmt.Errorf("bringing %q down: %w", name, err)
		}
	}
	return nil
}

func (r *LinkAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LinkAttachmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(&data); err != nil {
		resp.Diagnostics.AddError("Failed to attach interface", err.Error())
		return
	}
	data.ID = types.StringValue(data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LinkAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tgtHandle, tgtNs, tgtOwned, err := r.targetHandle(data.Netns.ValueString())
	if err != nil {
		// Target namespace gone (e.g. container removed): the attachment is gone.
		resp.State.RemoveResource(ctx)
		return
	}
	if tgtOwned {
		defer tgtNs.Close()
		if tgtHandle != r.client.Handle {
			defer tgtHandle.Close()
		}
	}

	link, err := tgtHandle.LinkByName(data.Name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	attrs := link.Attrs()
	if attrs.MasterIndex > 0 {
		if master, err := tgtHandle.LinkByIndex(attrs.MasterIndex); err == nil {
			data.Master = types.StringValue(master.Attrs().Name)
		}
	} else {
		data.Master = types.StringNull()
	}
	data.Up = types.BoolValue(attrs.Flags&net.FlagUp != 0)
	data.ID = types.StringValue(data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LinkAttachmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(&data); err != nil {
		resp.Diagnostics.AddError("Failed to update interface attachment", err.Error())
		return
	}
	data.ID = types.StringValue(data.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LinkAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LinkAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	netnsSpec := data.Netns.ValueString()
	tgtHandle, tgtNs, tgtOwned, err := r.targetHandle(netnsSpec)
	if err != nil {
		// Target namespace already gone: nothing to undo.
		return
	}
	if tgtOwned {
		defer tgtNs.Close()
		if tgtHandle != r.client.Handle {
			defer tgtHandle.Close()
		}
	}

	link, err := tgtHandle.LinkByName(data.Name.ValueString())
	if err != nil {
		return
	}

	// Detach from master and bring down.
	if link.Attrs().MasterIndex > 0 {
		_ = tgtHandle.LinkSetNoMaster(link)
	}
	_ = tgtHandle.LinkSetDown(link)

	// Move it back to the provider's namespace if we had moved it out.
	if netnsSpec != "" {
		srcNs, owned, err := r.sourceNs()
		if err == nil {
			if owned {
				defer srcNs.Close()
			}
			_ = tgtHandle.LinkSetNsFd(link, int(srcNs))
		}
	}
}

func (r *LinkAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
