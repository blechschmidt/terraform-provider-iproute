package provider_test

import (
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProvider_noNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "iproute" {}
data "iproute_link" "lo" {
  name = "lo"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_link.lo", "name", "lo"),
				),
			},
		},
	})
}

func TestAccProvider_invalidNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "iproute" {
  namespace = "nonexistent-ns-12345"
}
data "iproute_link" "lo" {
  name = "lo"
}
`,
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccProvider_emptyNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "iproute" {
  namespace = ""
}
data "iproute_link" "lo" {
  name = "lo"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.iproute_link.lo", "name", "lo"),
				),
			},
		},
	})
}
