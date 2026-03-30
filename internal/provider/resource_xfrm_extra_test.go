package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccXfrmStateResource_transport_crypt(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_state" "test" {
  src   = "10.0.0.1"
  dst   = "10.0.0.2"
  proto = "esp"
  spi   = 3000
  mode  = "transport"
  crypt {
    name = "cbc(aes)"
    key  = "0123456789abcdef0123456789abcdef"
  }
  auth {
    name = "hmac(sha256)"
    key  = "0123456789abcdef0123456789abcdef"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "mode", "transport"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "spi", "3000"),
				),
			},
		},
	})
}

func TestAccXfrmStateResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_state" "test" {
  src   = "fd00::1"
  dst   = "fd00::2"
  proto = "esp"
  spi   = 4000
  mode  = "tunnel"
  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "0123456789abcdef0123456789abcdef01234567"
    icv_len = 128
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "src", "fd00::1"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "dst", "fd00::2"),
				),
			},
		},
	})
}

func TestAccXfrmPolicyResource_forward(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_policy" "test" {
  src = "10.0.0.0/24"
  dst = "10.0.3.0/24"
  dir = "fwd"
  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dir", "fwd"),
				),
			},
		},
	})
}

func TestAccXfrmPolicyResource_mark(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_policy" "test" {
  src  = "10.0.0.0/24"
  dst  = "10.0.4.0/24"
  dir  = "out"
  mark = 42
  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "mark", "42"),
				),
			},
		},
	})
}

func TestAccXfrmPolicyResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_policy" "test" {
  src = "fd00::/64"
  dst = "fd01::/64"
  dir = "out"
  templates = [{
    src   = "fd00::1"
    dst   = "fd00::2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "src", "fd00::/64"),
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dst", "fd01::/64"),
				),
			},
		},
	})
}
