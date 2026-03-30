package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNexthopResource_ipv6Gateway(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-nhv60"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
resource "iproute_nexthop" "test" {
  nhid    = 300
  gateway = "fd00::254"
  device  = iproute_link.test.name
  family  = "inet6"
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_nexthop.test", "nhid", "300"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "gateway", "fd00::254"),
					resource.TestCheckResourceAttr("iproute_nexthop.test", "family", "inet6"),
				),
			},
		},
	})
}

func TestAccNexthopResource_multipleDevices(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "d1" {
  name = "test-nhmd1"
  type = "dummy"
  dummy {}
}
resource "iproute_link" "d2" {
  name = "test-nhmd2"
  type = "dummy"
  dummy {}
}
resource "iproute_nexthop" "n1" {
  nhid   = 301
  device = iproute_link.d1.name
}
resource "iproute_nexthop" "n2" {
  nhid   = 302
  device = iproute_link.d2.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_nexthop.n1", "nhid", "301"),
					resource.TestCheckResourceAttr("iproute_nexthop.n2", "nhid", "302"),
				),
			},
		},
	})
}
