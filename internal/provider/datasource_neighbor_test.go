package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNeighborDataSource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-dsneigh0"
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
data "iproute_neighbor" "test" {
  device     = iproute_link.test.name
  depends_on = [iproute_neighbor.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_neighbor.test", "neighbors.#"),
				),
			},
		},
	})
}

func TestAccNeighborDataSource_all(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
data "iproute_neighbor" "test" {}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_neighbor.test", "neighbors.#"),
				),
			},
		},
	})
}
