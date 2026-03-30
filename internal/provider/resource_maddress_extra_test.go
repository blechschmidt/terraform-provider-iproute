package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMaddressResource_multipleAddresses(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-maddr1"
  type = "dummy"
  dummy {}
}
resource "iproute_maddress" "test1" {
  device  = iproute_link.test.name
  address = "01:00:5e:00:00:01"
}
resource "iproute_maddress" "test2" {
  device  = iproute_link.test.name
  address = "01:00:5e:00:00:02"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_maddress.test1", "address", "01:00:5e:00:00:01"),
					resource.TestCheckResourceAttr("iproute_maddress.test2", "address", "01:00:5e:00:00:02"),
				),
			},
		},
	})
}
