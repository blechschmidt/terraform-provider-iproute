package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetnsResource_basic(t *testing.T) {
	nsName := fmt.Sprintf("tf-test-ns-%d", testutils.RandInt())

	t.Cleanup(func() {
		testutils.IPExec(nsName) // ignore cleanup errors
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetnsConfig(nsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_netns.test", "name", nsName),
				),
			},
			// Import
			{
				ResourceName:      "iproute_netns.test",
				ImportState:       true,
				ImportStateId:     nsName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetnsConfig(name string) string {
	return fmt.Sprintf(`
provider "iproute" {}
resource "iproute_netns" "test" {
  name = %q
}
`, name)
}
