package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MACAddress validates that a string is a valid MAC address.
type MACAddress struct{}

func (v MACAddress) Description(_ context.Context) string {
	return "value must be a valid MAC address"
}

func (v MACAddress) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v MACAddress) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	_, err := net.ParseMAC(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MAC Address",
			fmt.Sprintf("%q is not a valid MAC address: %s", value, err),
		)
	}
}

// IsMACAddress returns a validator that checks for valid MAC addresses.
func IsMACAddress() validator.String {
	return MACAddress{}
}
