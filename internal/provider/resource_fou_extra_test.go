package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFouResource_withProtocol(t *testing.T) {
	ensureFouModule(t)
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_fou" "test" {
  port       = 5560
  protocol   = 47
  encap_type = "direct"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_fou.test", "port", "5560"),
					resource.TestCheckResourceAttr("iproute_fou.test", "protocol", "47"),
				),
			},
		},
	})
}

func TestAccFouResource_localPort(t *testing.T) {
	ensureFouModule(t)
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_fou" "test" {
  port       = 5570
  encap_type = "gue"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_fou.test", "port", "5570"),
					resource.TestCheckResourceAttr("iproute_fou.test", "encap_type", "gue"),
				),
			},
		},
	})
}
