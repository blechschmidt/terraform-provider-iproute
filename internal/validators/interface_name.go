package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// InterfaceName validates that a string is a valid Linux interface name.
type InterfaceName struct{}

func (v InterfaceName) Description(_ context.Context) string {
	return "value must be a valid Linux interface name (1-15 characters, no spaces or /)"
}

func (v InterfaceName) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v InterfaceName) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if len(value) == 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interface Name",
			"Interface name must not be empty.",
		)
		return
	}

	if len(value) > 15 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interface Name",
			fmt.Sprintf("Interface name %q exceeds maximum length of 15 characters.", value),
		)
		return
	}

	if strings.Contains(value, "/") || strings.Contains(value, " ") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interface Name",
			fmt.Sprintf("Interface name %q must not contain '/' or spaces.", value),
		)
		return
	}

	if value == "." || value == ".." {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interface Name",
			fmt.Sprintf("Interface name %q is not allowed.", value),
		)
	}
}

// IsInterfaceName returns a validator that checks for valid Linux interface names.
func IsInterfaceName() validator.String {
	return InterfaceName{}
}
