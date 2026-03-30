package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_geneveWithPort(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-gnvp0"
  type = "geneve"
  geneve {
    vni    = 500
    remote = "10.0.0.2"
    port   = 6081
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-gnvp0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "geneve"),
				),
			},
		},
	})
}

func TestAccLinkResource_geneveMultiple(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "g1" {
  name = "test-gnvm1"
  type = "geneve"
  geneve {
    vni    = 600
    remote = "10.0.0.2"
  }
}
resource "iproute_link" "g2" {
  name = "test-gnvm2"
  type = "geneve"
  geneve {
    vni    = 601
    remote = "10.0.0.3"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.g1", "name", "test-gnvm1"),
					resource.TestCheckResourceAttr("iproute_link.g2", "name", "test-gnvm2"),
				),
			},
		},
	})
}
