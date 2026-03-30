package validators_test

import (
	"context"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsIPAddress_valid(t *testing.T) {
	tests := []string{
		"10.0.0.1",
		"192.168.1.1",
		"255.255.255.255",
		"0.0.0.0",
		"::1",
		"fd00::1",
		"2001:db8::1",
		"fe80::1",
	}

	v := validators.IsIPAddress()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be valid, got error: %s", tc, resp.Diagnostics.Errors()[0].Detail())
			}
		})
	}
}

func TestIsIPAddress_invalid(t *testing.T) {
	tests := []string{
		"not-an-ip",
		"10.0.0.256",
		"10.0.0",
		"",
		"abc",
		"10.0.0.1/24",
	}

	v := validators.IsIPAddress()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be invalid", tc)
			}
		})
	}
}

func TestIsIPAddress_null(t *testing.T) {
	v := validators.IsIPAddress()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("null value should not produce error")
	}
}

func TestIsIPAddress_unknown(t *testing.T) {
	v := validators.IsIPAddress()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringUnknown(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("unknown value should not produce error")
	}
}

func TestIsCIDR_valid(t *testing.T) {
	tests := []string{
		"10.0.0.0/24",
		"10.0.0.0/8",
		"192.168.1.0/32",
		"0.0.0.0/0",
		"fd00::/64",
		"::/0",
		"2001:db8::/32",
	}

	v := validators.IsCIDR()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be valid, got error: %s", tc, resp.Diagnostics.Errors()[0].Detail())
			}
		})
	}
}

func TestIsCIDR_invalid(t *testing.T) {
	tests := []string{
		"not-a-cidr",
		"10.0.0.1",
		"10.0.0.0/33",
		"",
		"abc/24",
	}

	v := validators.IsCIDR()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be invalid", tc)
			}
		})
	}
}

func TestIsCIDR_null(t *testing.T) {
	v := validators.IsCIDR()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("null value should not produce error")
	}
}

func TestIsMACAddress_valid(t *testing.T) {
	tests := []string{
		"aa:bb:cc:dd:ee:ff",
		"00:00:00:00:00:00",
		"ff:ff:ff:ff:ff:ff",
		"01:23:45:67:89:ab",
		"AA:BB:CC:DD:EE:FF",
	}

	v := validators.IsMACAddress()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be valid, got error: %s", tc, resp.Diagnostics.Errors()[0].Detail())
			}
		})
	}
}

func TestIsMACAddress_invalid(t *testing.T) {
	tests := []string{
		"not-a-mac",
		"aa:bb:cc:dd:ee",
		"aa:bb:cc:dd:ee:gg",
		"",
		"zz:zz:zz:zz:zz:zz",
	}

	v := validators.IsMACAddress()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be invalid", tc)
			}
		})
	}
}

func TestIsMACAddress_null(t *testing.T) {
	v := validators.IsMACAddress()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("null value should not produce error")
	}
}

func TestIsInterfaceName_valid(t *testing.T) {
	tests := []string{
		"eth0",
		"lo",
		"br-lan",
		"veth_1234",
		"wg0",
		"a",
		"123456789012345", // 15 chars
	}

	v := validators.IsInterfaceName()
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be valid, got error: %s", tc, resp.Diagnostics.Errors()[0].Detail())
			}
		})
	}
}

func TestIsInterfaceName_invalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"too_long", "1234567890123456"},
		{"has_slash", "eth/0"},
		{"has_space", "eth 0"},
		{"dot", "."},
		{"dotdot", ".."},
	}

	v := validators.IsInterfaceName()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tc.value),
			}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("expected %q to be invalid", tc.value)
			}
		})
	}
}

func TestIsInterfaceName_null(t *testing.T) {
	v := validators.IsInterfaceName()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("null value should not produce error")
	}
}

func TestIsInterfaceName_unknown(t *testing.T) {
	v := validators.IsInterfaceName()
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringUnknown(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Error("unknown value should not produce error")
	}
}

func TestIsIPAddress_description(t *testing.T) {
	v := validators.IsIPAddress()
	desc := v.Description(context.Background())
	if desc == "" {
		t.Error("description should not be empty")
	}
	md := v.MarkdownDescription(context.Background())
	if md == "" {
		t.Error("markdown description should not be empty")
	}
}

func TestIsCIDR_description(t *testing.T) {
	v := validators.IsCIDR()
	desc := v.Description(context.Background())
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestIsMACAddress_description(t *testing.T) {
	v := validators.IsMACAddress()
	desc := v.Description(context.Background())
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestIsInterfaceName_description(t *testing.T) {
	v := validators.IsInterfaceName()
	desc := v.Description(context.Background())
	if desc == "" {
		t.Error("description should not be empty")
	}
}
