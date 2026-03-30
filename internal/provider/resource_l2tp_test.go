package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2tpTunnelResource_udp(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_l2tp_tunnel" "test" {
  tunnel_id      = 100
  peer_tunnel_id = 200
  encap_type     = "udp"
  local          = "127.0.0.1"
  remote         = "127.0.0.2"
  local_port     = 5000
  remote_port    = 5001
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_l2tp_tunnel.test", "tunnel_id", "100"),
					resource.TestCheckResourceAttr("iproute_l2tp_tunnel.test", "peer_tunnel_id", "200"),
					resource.TestCheckResourceAttr("iproute_l2tp_tunnel.test", "encap_type", "udp"),
				),
			},
		},
	})
}

func TestAccL2tpTunnelResource_ip(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_l2tp_tunnel" "test" {
  tunnel_id      = 101
  peer_tunnel_id = 201
  encap_type     = "ip"
  local          = "127.0.0.1"
  remote         = "127.0.0.2"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_l2tp_tunnel.test", "tunnel_id", "101"),
					resource.TestCheckResourceAttr("iproute_l2tp_tunnel.test", "encap_type", "ip"),
				),
			},
		},
	})
}

func TestAccL2tpSessionResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_l2tp_tunnel" "test" {
  tunnel_id      = 110
  peer_tunnel_id = 210
  encap_type     = "udp"
  local          = "127.0.0.1"
  remote         = "127.0.0.2"
  local_port     = 6000
  remote_port    = 6001
}
resource "iproute_l2tp_session" "test" {
  tunnel_id       = iproute_l2tp_tunnel.test.tunnel_id
  session_id      = 1000
  peer_session_id = 2000
  name            = "l2tp-sess0"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_l2tp_session.test", "session_id", "1000"),
					resource.TestCheckResourceAttr("iproute_l2tp_session.test", "peer_session_id", "2000"),
				),
			},
		},
	})
}
