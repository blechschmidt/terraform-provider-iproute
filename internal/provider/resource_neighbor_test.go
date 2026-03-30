package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNeighborResource_ipv4(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNeighborConfig_ipv4(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_neighbor.test", "address", "10.0.0.2"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "lladdr", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "family", "inet"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "state", "permanent"),
				),
			},
			// Import
			{
				ResourceName:      "iproute_neighbor.test",
				ImportState:       true,
				ImportStateId:     "test-neigh0|10.0.0.2",
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccNeighborResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNeighborConfig_ipv6(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_neighbor.test", "address", "fd00::2"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "family", "inet6"),
				),
			},
		},
	})
}

func TestAccNeighborResource_update_lladdr(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNeighborConfig_lladdr(ns, "aa:bb:cc:dd:ee:ff"),
				Check:  resource.TestCheckResourceAttr("iproute_neighbor.test", "lladdr", "aa:bb:cc:dd:ee:ff"),
			},
			{
				Config: testAccNeighborConfig_lladdr(ns, "11:22:33:44:55:66"),
				Check:  resource.TestCheckResourceAttr("iproute_neighbor.test", "lladdr", "11:22:33:44:55:66"),
			},
		},
	})
}

func testAccNeighborConfig_ipv4(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-neigh0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address = "10.0.0.2"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = iproute_link.test.name
  state   = "permanent"
  depends_on = [iproute_address.test]
}
`, ns)
}

func testAccNeighborConfig_ipv6(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-neigh6"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address = "fd00::2"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = iproute_link.test.name
  state   = "permanent"
  depends_on = [iproute_address.test]
}
`, ns)
}

func testAccNeighborConfig_lladdr(ns, mac string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-nupd0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address = "10.0.0.2"
  lladdr  = %q
  device  = iproute_link.test.name
  state   = "permanent"
  depends_on = [iproute_address.test]
}
`, ns, mac)
}
