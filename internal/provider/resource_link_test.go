package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_dummy(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_dummy(ns, "test-dummy0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-dummy0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "dummy"),
					resource.TestCheckResourceAttr("iproute_link.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("iproute_link.test", "if_index"),
					resource.TestCheckResourceAttrSet("iproute_link.test", "mac_address"),
				),
			},
			// Update: disable
			{
				Config: testAccLinkConfig_dummyDisabled(ns, "test-dummy0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "enabled", "false"),
					resource.TestCheckResourceAttr("iproute_link.test", "admin_status", "down"),
				),
			},
			// Import
			{
				ResourceName:      "iproute_link.test",
				ImportState:       true,
				ImportStateId:     "test-dummy0",
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccLinkResource_bridge(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_bridge(ns, "test-br0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-br0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "bridge"),
				),
			},
		},
	})
}

func TestAccLinkResource_veth(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_veth(ns, "test-veth0", "test-veth1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-veth0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "veth"),
				),
			},
		},
	})
}

func TestAccLinkResource_vxlan(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_vxlan(ns, "test-vxlan0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-vxlan0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "vxlan"),
				),
			},
		},
	})
}

func TestAccLinkResource_bond(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_bond(ns, "test-bond0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-bond0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "bond"),
				),
			},
		},
	})
}

func TestAccLinkResource_wireguard(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_wireguard(ns, "test-wg0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-wg0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "wireguard"),
				),
			},
		},
	})
}

func TestAccLinkResource_gre(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_gre(ns, "test-gre0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-gre0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "gre"),
				),
			},
		},
	})
}

func TestAccLinkResource_ipip(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_ipip(ns, "test-ipip0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-ipip0"),
				),
			},
		},
	})
}

func TestAccLinkResource_sit(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_sit(ns, "test-sit0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-sit0"),
				),
			},
		},
	})
}

func TestAccLinkResource_geneve(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_geneve(ns, "test-gnv0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-gnv0"),
				),
			},
		},
	})
}

func TestAccLinkResource_tuntap(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_tuntap(ns, "test-tun0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-tun0"),
				),
			},
		},
	})
}

func TestAccLinkResource_mtu_update(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLinkConfig_dummyWithMTU(ns, "test-mtu0", 1400),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "mtu", "1400"),
				),
			},
			{
				Config: testAccLinkConfig_dummyWithMTU(ns, "test-mtu0", 9000),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "mtu", "9000"),
				),
			},
		},
	})
}

// --- Config helpers ---

func testAccLinkConfig_dummy(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "dummy"
  dummy {}
}
`, ns, name)
}

func testAccLinkConfig_dummyDisabled(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name    = %q
  type    = "dummy"
  enabled = false
  dummy {}
}
`, ns, name)
}

func testAccLinkConfig_dummyWithMTU(ns, name string, mtu int) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "dummy"
  mtu  = %d
  dummy {}
}
`, ns, name, mtu)
}

func testAccLinkConfig_bridge(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "bridge"
  bridge {}
}
`, ns, name)
}

func testAccLinkConfig_veth(ns, name, peer string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "veth"
  veth {
    peer_name = %q
  }
}
`, ns, name, peer)
}

func testAccLinkConfig_vxlan(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "vxlan"
  vxlan {
    vni  = 100
    port = 4789
  }
}
`, ns, name)
}

func testAccLinkConfig_bond(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "bond"
  bond {
    mode   = "balance-rr"
    miimon = 100
  }
}
`, ns, name)
}

func testAccLinkConfig_wireguard(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "wireguard"
  wireguard {}
}
`, ns, name)
}

func testAccLinkConfig_gre(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "gre"
  gre {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}
`, ns, name)
}

func testAccLinkConfig_ipip(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "ipip"
  ipip {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}
`, ns, name)
}

func testAccLinkConfig_sit(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "sit"
  sit {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}
`, ns, name)
}

func testAccLinkConfig_geneve(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "geneve"
  geneve {
    vni    = 200
    remote = "10.0.0.2"
  }
}
`, ns, name)
}

func testAccLinkConfig_tuntap(ns, name string) string {
	return fmt.Sprintf(`
provider "iproute" {
  namespace = %q
}
resource "iproute_link" "test" {
  name = %q
  type = "tuntap"
  tuntap {
    mode = "tun"
  }
}
`, ns, name)
}
