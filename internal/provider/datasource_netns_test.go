package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetnsDataSource_basic(t *testing.T) {
	nsName := fmt.Sprintf("tf-test-dsns-%d", testutils.RandInt())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" {}
resource "iproute_netns" "test" {
  name = %q
}
data "iproute_netns" "test" {
  depends_on = [iproute_netns.test]
}
`, nsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iproute_netns.test", "namespaces.#"),
				),
			},
		},
	})
}
