package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRuleResource_dst(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 2000
  dst      = "10.1.0.0/24"
  table    = 200
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "priority", "2000"),
					resource.TestCheckResourceAttr("iproute_rule.test", "dst", "10.1.0.0/24"),
					resource.TestCheckResourceAttr("iproute_rule.test", "table", "200"),
				),
			},
		},
	})
}

func TestAccRuleResource_iif(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-ruliif0"
  type = "dummy"
  dummy {}
}
resource "iproute_rule" "test" {
  priority = 3000
  iif_name = iproute_link.test.name
  table    = 300
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "iif_name", "test-ruliif0"),
					resource.TestCheckResourceAttr("iproute_rule.test", "table", "300"),
				),
			},
		},
	})
}

func TestAccRuleResource_oif(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-ruloif0"
  type = "dummy"
  dummy {}
}
resource "iproute_rule" "test" {
  priority = 3001
  oif_name = iproute_link.test.name
  table    = 301
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "oif_name", "test-ruloif0"),
				),
			},
		},
	})
}

func TestAccRuleResource_tos(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 4000
  src      = "10.0.0.0/24"
  tos      = 16
  table    = 400
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "tos", "16"),
				),
			},
		},
	})
}

func TestAccRuleResource_ipv6(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 5000
  src      = "fd00::/64"
  table    = 500
  family   = "inet6"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "src", "fd00::/64"),
					resource.TestCheckResourceAttr("iproute_rule.test", "family", "inet6"),
				),
			},
		},
	})
}

func TestAccRuleResource_update(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 6000
  src      = "10.0.0.0/24"
  table    = 600
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_rule.test", "table", "600"),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "test" {
  priority = 6000
  src      = "10.0.0.0/24"
  table    = 601
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_rule.test", "table", "601"),
			},
		},
	})
}
