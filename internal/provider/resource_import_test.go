package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteResource_import(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-imprt0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.1.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "10.1.0.0/24"),
				),
			},
			{
				ResourceName:      "iproute_route.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccRuleResource_import(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 1000
  src      = "10.0.0.0/24"
  table    = 100
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "priority", "1000"),
				),
			},
			{
				ResourceName:      "iproute_rule.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccNeighborResource_import(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-impne0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address    = "10.0.0.2"
  lladdr     = "aa:bb:cc:dd:ee:ff"
  device     = iproute_link.test.name
  state      = "permanent"
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_neighbor.test", "address", "10.0.0.2"),
				),
			},
			{
				ResourceName:      "iproute_neighbor.test",
				ImportState:       true,
				ImportStateId:     "test-impne0|10.0.0.2",
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccLinkResource_import_bridge(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-impbr0"
  type = "bridge"
  bridge {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-impbr0"),
				),
			},
			{
				ResourceName:      "iproute_link.test",
				ImportState:       true,
				ImportStateId:     "test-impbr0",
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAddressResource_import_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-impa60"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "address", "fd00::1/64"),
				),
			},
			{
				ResourceName:      "iproute_address.test",
				ImportState:       true,
				ImportStateId:     "test-impa60|fd00::1/64",
				ImportStateVerify: false,
			},
		},
	})
}
