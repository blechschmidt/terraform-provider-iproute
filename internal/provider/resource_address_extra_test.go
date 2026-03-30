package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAddressResource_scope(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-addrsc0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
  scope   = "host"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "scope", "host"),
				),
			},
		},
	})
}

func TestAccAddressResource_broadcast(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-addrbc0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address   = "10.0.0.1/24"
  device    = iproute_link.test.name
  broadcast = "10.0.0.255"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "broadcast", "10.0.0.255"),
				),
			},
		},
	})
}

func TestAccAddressResource_peer(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-addrpe0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/32"
  device  = iproute_link.test.name
  peer    = "10.0.0.2/32"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.test", "peer", "10.0.0.2/32"),
				),
			},
		},
	})
}

func TestAccAddressResource_multipleIpv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-addrmv6"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "a1" {
  address = "fd00::1/64"
  device  = iproute_link.test.name
}
resource "iproute_address" "a2" {
  address = "fd00::2/64"
  device  = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.a1", "address", "fd00::1/64"),
					resource.TestCheckResourceAttr("iproute_address.a2", "address", "fd00::2/64"),
				),
			},
		},
	})
}
