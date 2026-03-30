package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNexthopDataSource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-dsnh0"
  type = "dummy"
  dummy {}
}
resource "iproute_nexthop" "test" {
  nhid   = 200
  device = iproute_link.test.name
}
data "iproute_nexthop" "test" {
  depends_on = [iproute_nexthop.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_nexthop.test", "nexthops.#"),
				),
			},
		},
	})
}
