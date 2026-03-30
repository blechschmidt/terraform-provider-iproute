package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAddressResource_ipv4(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddressConfig_ipv4(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "address", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("iproute_address.test", "device", "test-addr0"),
					resource.TestCheckResourceAttr("iproute_address.test", "family", "inet"),
				),
			},
			// Import
			{
				ResourceName:      "iproute_address.test",
				ImportState:       true,
				ImportStateId:     "test-addr0|10.0.0.1/24",
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAddressResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddressConfig_ipv6(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "address", "fd00::1/64"),
					resource.TestCheckResourceAttr("iproute_address.test", "family", "inet6"),
				),
			},
		},
	})
}

func TestAccAddressResource_label(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddressConfig_label(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "label", "test-addr0:web"),
				),
			},
		},
	})
}

func TestAccAddressResource_multiAddr(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddressConfig_multi(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test1", "address", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("iproute_address.test2", "address", "10.0.0.2/24"),
				),
			},
		},
	})
}

func testAccAddressConfig_ipv4(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-addr0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
`, ns)
}

func testAccAddressConfig_ipv6(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-addr6"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
`, ns)
}

func testAccAddressConfig_label(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-addr0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
  label   = "test-addr0:web"
}
`, ns)
}

func testAccAddressConfig_multi(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = "test-multi0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test1" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_address" "test2" {
  address = "10.0.0.2/24"
  device  = iproute_link.test.name
}
`, ns)
}
