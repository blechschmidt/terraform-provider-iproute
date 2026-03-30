package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRuleResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleConfig_basic(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "src", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("iproute_rule.test", "table", "100"),
					resource.TestCheckResourceAttr("iproute_rule.test", "priority", "100"),
				),
			},
		},
	})
}

func TestAccRuleResource_fwmark(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleConfig_fwmark(ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "fwmark", "42"),
					resource.TestCheckResourceAttr("iproute_rule.test", "table", "200"),
				),
			},
		},
	})
}

func testAccRuleConfig_basic(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_rule" "test" {
  src      = "10.0.0.0/24"
  table    = 100
  priority = 100
}
`, ns)
}

func testAccRuleConfig_fwmark(ns string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_rule" "test" {
  fwmark   = 42
  table    = 200
  priority = 200
}
`, ns)
}
