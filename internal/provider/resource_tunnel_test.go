package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelResource_gre(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tgre0"
  mode   = "gre"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tunnel.test", "name", "test-tgre0"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "mode", "gre"),
				),
			},
		},
	})
}

func TestAccTunnelResource_ipip(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tipip0"
  mode   = "ipip"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tunnel.test", "name", "test-tipip0"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "mode", "ipip"),
				),
			},
		},
	})
}

func TestAccTunnelResource_sit(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tsit0"
  mode   = "sit"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tunnel.test", "name", "test-tsit0"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "mode", "sit"),
				),
			},
		},
	})
}
