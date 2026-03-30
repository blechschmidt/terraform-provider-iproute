package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteResource_unreachable(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_route" "test" {
  destination = "10.99.0.0/24"
  type        = "unreachable"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "10.99.0.0/24"),
					resource.TestCheckResourceAttr("iproute_route.test", "type", "unreachable"),
				),
			},
		},
	})
}

func TestAccRouteResource_prohibit(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_route" "test" {
  destination = "10.98.0.0/24"
  type        = "prohibit"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "10.98.0.0/24"),
					resource.TestCheckResourceAttr("iproute_route.test", "type", "prohibit"),
				),
			},
		},
	})
}

func TestAccRouteResource_scope(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rtscope"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.3.0.0/24"
  device      = iproute_link.test.name
  scope       = "link"
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "scope", "link"),
				),
			},
		},
	})
}

func TestAccRouteResource_protocol(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rtproto"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.4.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  protocol    = "static"
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "protocol", "static"),
				),
			},
		},
	})
}

func TestAccRouteResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rtv6"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "fd01::/64"
  gateway     = "fd00::2"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "fd01::/64"),
				),
			},
		},
	})
}

func TestAccRouteResource_source(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rtsrc"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.5.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  source      = "10.0.0.1"
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "source", "10.0.0.1"),
				),
			},
		},
	})
}

func TestAccRouteResource_mtu(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rtmtu"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.6.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  mtu         = 1400
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.test", "mtu", "1400"),
				),
			},
		},
	})
}
