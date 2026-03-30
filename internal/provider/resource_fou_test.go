package provider_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ensureFouModule(t *testing.T) {
	t.Helper()
	if err := exec.Command("modprobe", "fou").Run(); err != nil {
		t.Skip("FOU kernel module not available, skipping")
	}
}

func TestAccFouResource_direct(t *testing.T) {
	ensureFouModule(t)
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_fou" "test" {
  port       = 5555
  protocol   = 4
  encap_type = "direct"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_fou.test", "port", "5555"),
					resource.TestCheckResourceAttr("iproute_fou.test", "encap_type", "direct"),
				),
			},
		},
	})
}

func TestAccFouResource_gue(t *testing.T) {
	ensureFouModule(t)
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_fou" "test" {
  port       = 5556
  encap_type = "gue"
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_fou.test", "port", "5556"),
					resource.TestCheckResourceAttr("iproute_fou.test", "encap_type", "gue"),
				),
			},
		},
	})
}
