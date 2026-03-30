package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNexthopResource_gateway(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-nh0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_nexthop" "test" {
  nhid    = 100
  gateway = "10.0.0.254"
  device  = iproute_link.test.name
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_nexthop.test", "nhid", "100"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "gateway", "10.0.0.254"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "family", "inet"),
				),
			},
		},
	})
}

func TestAccNexthopResource_device(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-nh1"
  type = "dummy"
  dummy {}
}
resource "iproute_nexthop" "test" {
  nhid   = 101
  device = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_nexthop.test", "nhid", "101"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "device", "test-nh1"),
				),
			},
		},
	})
}

func TestAccNexthopResource_blackhole(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_nexthop" "test" {
  nhid      = 102
  blackhole = true
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_nexthop.test", "nhid", "102"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "blackhole", "true"),
				),
			},
		},
	})
}
