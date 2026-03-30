package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetnsResource_multiple(t *testing.T) {
	ns1 := fmt.Sprintf("tf-test-ns1-%d", testutils.RandInt())
	ns2 := fmt.Sprintf("tf-test-ns2-%d", testutils.RandInt())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" {}
resource "iproute_netns" "ns1" {
  name = %q
}
resource "iproute_netns" "ns2" {
  name = %q
}
`, ns1, ns2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_netns.ns1", "name", ns1),
					resource.TestCheckResourceAttr("iproute_netns.ns2", "name", ns2),
				),
			},
		},
	})
}

func TestAccNetnsResource_recreate(t *testing.T) {
	ns := fmt.Sprintf("tf-test-nsrc-%d", testutils.RandInt())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" {}
resource "iproute_netns" "test" {
  name = %q
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_netns.test", "name", ns),
			},
			// Destroy and recreate
			{
				Config: `provider "iproute" {}`,
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" {}
resource "iproute_netns" "test" {
  name = %q
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_netns.test", "name", ns),
			},
		},
	})
}
