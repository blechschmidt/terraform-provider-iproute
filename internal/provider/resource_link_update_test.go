package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_enableDisable(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-enadis0"
  type    = "dummy"
  enabled = true
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "enabled", "true"),
					resource.TestCheckResourceAttr("iproute_link.test", "admin_status", "up"),
				),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-enadis0"
  type    = "dummy"
  enabled = false
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "enabled", "false"),
					resource.TestCheckResourceAttr("iproute_link.test", "admin_status", "down"),
				),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name    = "test-enadis0"
  type    = "dummy"
  enabled = true
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "enabled", "true"),
					resource.TestCheckResourceAttr("iproute_link.test", "admin_status", "up"),
				),
			},
		},
	})
}

func TestAccLinkResource_updateMacAddress(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-macupd0"
  type        = "dummy"
  mac_address = "aa:bb:cc:dd:ee:01"
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "mac_address", "aa:bb:cc:dd:ee:01"),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-macupd0"
  type        = "dummy"
  mac_address = "aa:bb:cc:dd:ee:02"
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "mac_address", "aa:bb:cc:dd:ee:02"),
			},
		},
	})
}

func TestAccLinkResource_updateDescription(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-descupd0"
  type        = "dummy"
  description = "first"
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "description", "first"),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-descupd0"
  type        = "dummy"
  description = "second"
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "description", "second"),
			},
		},
	})
}

func TestAccLinkResource_updateTxQueueLen(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name         = "test-txqupd0"
  type         = "dummy"
  tx_queue_len = 500
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "tx_queue_len", "500"),
			},
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name         = "test-txqupd0"
  type         = "dummy"
  tx_queue_len = 1000
  dummy {}
}
`, ns),
				Check: resource.TestCheckResourceAttr("iproute_link.test", "tx_queue_len", "1000"),
			},
		},
	})
}

func TestAccLinkResource_bondModes(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-bond802"
  type = "bond"
  bond {
    mode   = "802.3ad"
    miimon = 100
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-bond802"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "bond"),
				),
			},
		},
	})
}

func TestAccLinkResource_vxlanExtended(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-vxlan1"
  type = "vxlan"
  vxlan {
    vni    = 300
    port   = 4789
    local  = "10.0.0.1"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-vxlan1"),
				),
			},
		},
	})
}
