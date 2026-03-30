package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTokenResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-token0"
  type    = "veth"
  enabled = true
  veth {
    peer_name = "test-token0p"
  }
}
resource "iproute_token" "test" {
  device = iproute_link.test.name
  token  = "::1"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_token.test", "device", "test-token0"),
					resource.TestCheckResourceAttr("iproute_token.test", "token", "::1"),
				),
			},
		},
	})
}

func TestAccTokenResource_update(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-token1"
  type    = "veth"
  enabled = true
  veth {
    peer_name = "test-token1p"
  }
}
resource "iproute_token" "test" {
  device = iproute_link.test.name
  token  = "::1"
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_token.test", "token", "::1"),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-token1"
  type    = "veth"
  enabled = true
  veth {
    peer_name = "test-token1p"
  }
}
resource "iproute_token" "test" {
  device = iproute_link.test.name
  token  = "::2"
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_token.test", "token", "::2"),
			},
		},
	})
}
