package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTcpMetricsResource_basic(t *testing.T) {
	t.Skip("TCP metrics entries are only created by the kernel via actual TCP connections")
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tcp_metrics" "test" {
  address = "10.0.0.1"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tcp_metrics.test", "address", "10.0.0.1"),
					resource.TestCheckResourceAttrSet("iproute_tcp_metrics.test", "rtt"),
					resource.TestCheckResourceAttrSet("iproute_tcp_metrics.test", "cwnd"),
				),
			},
		},
	})
}

func TestAccTcpMetricsResource_ipv6(t *testing.T) {
	t.Skip("TCP metrics entries are only created by the kernel via actual TCP connections")
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tcp_metrics" "test" {
  address = "fd00::1"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tcp_metrics.test", "address", "fd00::1"),
				),
			},
		},
	})
}
