package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelResource_greWithTTL(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tgrttl0"
  mode   = "gre"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tunnel.test", "name", "test-tgrttl0"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "mode", "gre"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "local", "10.0.0.1"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "remote", "10.0.0.2"),
				),
			},
		},
	})
}

func TestAccTunnelResource_ipipLocal(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tipipl0"
  mode   = "ipip"
  local  = "10.0.0.3"
  remote = "10.0.0.4"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tunnel.test", "local", "10.0.0.3"),
					resource.TestCheckResourceAttr("iproute_tunnel.test", "remote", "10.0.0.4"),
				),
			},
		},
	})
}
