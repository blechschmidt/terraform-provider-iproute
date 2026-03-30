package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// CIDR validates that a string is a valid CIDR notation.
type CIDR struct{}

func (v CIDR) Description(_ context.Context) string {
	return "value must be a valid CIDR notation (e.g., 10.0.0.1/24)"
}

func (v CIDR) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v CIDR) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	_, _, err := net.ParseCIDR(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			fmt.Sprintf("%q is not valid CIDR notation: %s", value, err),
		)
	}
}

// IsCIDR returns a validator that checks for valid CIDR notation.
func IsCIDR() validator.String {
	return CIDR{}
}
