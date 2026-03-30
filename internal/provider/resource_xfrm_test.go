package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccXfrmStateResource_esp(t *testing.T) {
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
  spi   = 1000
  mode  = "tunnel"
  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "0123456789abcdef0123456789abcdef01234567"
    icv_len = 128
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "src", "10.0.0.1"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "dst", "10.0.0.2"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "proto", "esp"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "spi", "1000"),
				),
			},
		},
	})
}

func TestAccXfrmPolicyResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_policy" "test" {
  src = "10.0.0.0/24"
  dst = "10.0.1.0/24"
  dir = "out"
  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "src", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dst", "10.0.1.0/24"),
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dir", "out"),
				),
			},
		},
	})
}

func TestAccXfrmStateResource_ah(t *testing.T) {
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
  proto = "ah"
  spi   = 2000
  mode  = "transport"
  auth {
    name = "hmac(sha256)"
    key  = "0123456789abcdef0123456789abcdef"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "proto", "ah"),
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "mode", "transport"),
				),
			},
		},
	})
}

func TestAccXfrmPolicyResource_block(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_policy" "test" {
  src    = "10.0.0.0/24"
  dst    = "10.0.2.0/24"
  dir    = "in"
  action = "block"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "action", "block"),
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dir", "in"),
				),
			},
		},
	})
}
