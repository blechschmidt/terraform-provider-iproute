package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_wireguardWithMTU(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-wgm0"
  type = "wireguard"
  mtu  = 1420
  wireguard {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-wgm0"),
					resource.TestCheckResourceAttr("iproute_link.test", "mtu", "1420"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "wireguard"),
				),
			},
		},
	})
}

func TestAccLinkResource_wireguardWithAddress(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-wga0"
  type = "wireguard"
  wireguard {}
}
resource "iproute_address" "test" {
  address = "10.99.0.1/24"
  device  = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "type", "wireguard"),
					resource.TestCheckResourceAttr("iproute_address.test", "address", "10.99.0.1/24"),
				),
			},
		},
	})
}
