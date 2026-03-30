package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_basic(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("iproute_route.test", "family", "inet"),
				),
			},
		},
	})
}

func TestAccRouteResource_blackhole(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_blackhole(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "type", "blackhole"),
				),
			},
		},
	})
}

func TestAccRouteResource_metric(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_metric(ns, 100),
				Check:  resource.TestCheckResourceAttr("iproute_route.test", "metric", "100"),
			},
			{
				Config: testAccRouteConfig_metric(ns, 200),
				Check:  resource.TestCheckResourceAttr("iproute_route.test", "metric", "200"),
			},
		},
	})
}

func TestAccRouteResource_table(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_table(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "table", "100"),
				),
			},
		},
	})
}

func testAccRouteConfig_basic(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-rte0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "192.168.1.0/24"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
`, ns)
}

func testAccRouteConfig_blackhole(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_route" "test" {
  destination = "192.168.99.0/24"
  type        = "blackhole"
}
`, ns)
}

func testAccRouteConfig_metric(ns string, metric int) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-metric0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "192.168.1.0/24"
  device      = iproute_link.test.name
  metric      = %d
  depends_on  = [iproute_address.test]
}
`, ns, metric)
}

func testAccRouteConfig_table(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-tbl0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "192.168.1.0/24"
  device      = iproute_link.test.name
  table       = 100
  depends_on  = [iproute_address.test]
}
`, ns)
}
