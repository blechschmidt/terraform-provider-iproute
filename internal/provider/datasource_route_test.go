package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteDataSource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-dsrt0"
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
data "iproute_route" "test" {
  family     = "inet"
  depends_on = [iproute_route.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.#"),
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.0.destination"),
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.0.family"),
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.0.protocol"),
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.0.scope"),
				),
			},
		},
	})
}

func TestAccRouteDataSource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
data "iproute_route" "test" {
  family = "inet6"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_route.test", "routes.#"),
				),
			},
		},
	})
}

func TestAccRouteDataSource_get(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-dsrtg0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "192.0.2.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "198.51.100.0/24"
  gateway     = "192.0.2.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
data "iproute_route" "test" {
  get        = "198.51.100.42"
  depends_on = [iproute_route.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_route.test", "routes.#", "1"),
					resource.TestCheckResourceAttr("data.iproute_route.test", "routes.0.gateway", "192.0.2.254"),
					resource.TestCheckResourceAttr("data.iproute_route.test", "routes.0.device", "test-dsrtg0"),
					resource.TestCheckResourceAttr("data.iproute_route.test", "routes.0.family", "inet"),
					resource.TestCheckResourceAttr("data.iproute_route.test", "id", "get:198.51.100.42"),
				),
			},
		},
	})
}
