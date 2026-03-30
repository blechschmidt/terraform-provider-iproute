package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNeighborResource_permanent(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-neighp0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address    = "10.0.0.100"
  lladdr     = "11:22:33:44:55:66"
  device     = iproute_link.test.name
  state      = "permanent"
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_neighbor.test", "state", "permanent"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "lladdr", "11:22:33:44:55:66"),
				),
			},
		},
	})
}

func TestAccNeighborResource_multiple(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-neighm0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "n1" {
  address    = "10.0.0.2"
  lladdr     = "aa:bb:cc:dd:ee:01"
  device     = iproute_link.test.name
  state      = "permanent"
  depends_on = [iproute_address.test]
}
resource "iproute_neighbor" "n2" {
  address    = "10.0.0.3"
  lladdr     = "aa:bb:cc:dd:ee:02"
  device     = iproute_link.test.name
  state      = "permanent"
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_neighbor.n1", "address", "10.0.0.2"),
					resource.TestCheckResourceAttr("iproute_neighbor.n2", "address", "10.0.0.3"),
				),
			},
		},
	})
}
