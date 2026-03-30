package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_bridgeVlanFiltering(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-brvf0"
  type = "bridge"
  bridge {
    vlan_filtering = true
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-brvf0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "bridge"),
				),
			},
		},
	})
}

func TestAccLinkResource_bridgeHelloTime(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-brht0"
  type = "bridge"
  bridge {
    hello_time = 200
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-brht0"),
				),
			},
		},
	})
}

func TestAccLinkResource_bridgeMultiplePorts(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "br" {
  name = "test-brmp0"
  type = "bridge"
  bridge {}
}
resource "iproute_link" "p1" {
  name   = "test-brmp1"
  type   = "dummy"
  master = iproute_link.br.name
  dummy {}
}
resource "iproute_link" "p2" {
  name   = "test-brmp2"
  type   = "dummy"
  master = iproute_link.br.name
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.p1", "master", "test-brmp0"),
					resource.TestCheckResourceAttr("iproute_link.p2", "master", "test-brmp0"),
				),
			},
		},
	})
}
