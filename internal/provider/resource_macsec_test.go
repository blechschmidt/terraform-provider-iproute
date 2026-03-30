package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMacsecResource_basic(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-macsec0"
  type = "dummy"
  dummy {}
}
resource "iproute_macsec" "test" {
  parent  = iproute_link.test.name
  name    = "macsec0"
  encrypt = true
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_macsec.test", "name", "macsec0"),
					resource.TestCheckResourceAttr("iproute_macsec.test", "parent", "test-macsec0"),
					resource.TestCheckResourceAttr("iproute_macsec.test", "encrypt", "true"),
				),
			},
		},
	})
}

func TestAccMacsecResource_withPort(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-macsec1"
  type = "dummy"
  dummy {}
}
resource "iproute_macsec" "test" {
  parent  = iproute_link.test.name
  name    = "macsec1"
  port    = 1
  encrypt = false
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_macsec.test", "name", "macsec1"),
					resource.TestCheckResourceAttr("iproute_macsec.test", "port", "1"),
				),
			},
		},
	})
}
