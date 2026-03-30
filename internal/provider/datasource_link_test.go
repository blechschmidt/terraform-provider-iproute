package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkDataSource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-dslink0"
  type = "dummy"
  mtu  = 1400
  dummy {}
}
data "iproute_link" "test" {
  name = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_link.test", "name", "test-dslink0"),
					resource.TestCheckResourceAttr("data.iproute_link.test", "type", "dummy"),
					resource.TestCheckResourceAttr("data.iproute_link.test", "mtu", "1400"),
					resource.TestCheckResourceAttrSet("data.iproute_link.test", "if_index"),
					resource.TestCheckResourceAttrSet("data.iproute_link.test", "mac_address"),
				),
			},
		},
	})
}

func TestAccLinkDataSource_loopback(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
data "iproute_link" "lo" {
  name = "lo"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_link.lo", "name", "lo"),
					resource.TestCheckResourceAttr("data.iproute_link.lo", "mtu", "65536"),
					resource.TestCheckResourceAttr("data.iproute_link.lo", "admin_status", "up"),
				),
			},
		},
	})
}

func TestAccLinkDataSource_bridge_master(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "br" {
  name = "test-dsbr0"
  type = "bridge"
  bridge {}
}
resource "iproute_link" "port" {
  name   = "test-dsport0"
  type   = "dummy"
  master = iproute_link.br.name
  dummy {}
}
data "iproute_link" "port" {
  name = iproute_link.port.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_link.port", "name", "test-dsport0"),
					resource.TestCheckResourceAttr("data.iproute_link.port", "master", "test-dsbr0"),
				),
			},
		},
	})
}
