package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTuntapResource_tunImport(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tuntap" "test" {
  name = "test-tunni0"
  mode = "tun"
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_tuntap.test", "mode", "tun"),
			},
		},
	})
}

func TestAccTuntapResource_multiQueue(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tuntap" "test" {
  name = "test-tunmq0"
  mode = "tap"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_tuntap.test", "name", "test-tunmq0"),
					resource.TestCheckResourceAttr("iproute_tuntap.test", "mode", "tap"),
				),
			},
		},
	})
}
