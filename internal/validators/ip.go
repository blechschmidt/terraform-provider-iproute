package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IPAddress validates that a string is a valid IP address.
type IPAddress struct{}

func (v IPAddress) Description(_ context.Context) string {
	return "value must be a valid IP address"
}

func (v IPAddress) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v IPAddress) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if net.ParseIP(value) == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid IP Address",
			fmt.Sprintf("%q is not a valid IP address.", value),
		)
	}
}

// IsIPAddress returns a validator that checks for valid IP addresses.
func IsIPAddress() validator.String {
	return IPAddress{}
}
