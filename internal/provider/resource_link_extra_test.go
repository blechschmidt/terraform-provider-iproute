package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Additional link tests: macvlan, ipvlan, vlan, vti, ip6tnl, ifb, macvtap

func TestAccLinkResource_macvlan(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "parent" {
  name = "test-mvparent"
  type = "dummy"
  dummy {}
}
resource "iproute_link" "test" {
  name   = "test-macvlan0"
  type   = "macvlan"
  macvlan {
    parent = iproute_link.parent.name
    mode   = "bridge"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-macvlan0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "macvlan"),
				),
			},
		},
	})
}

func TestAccLinkResource_ipvlan(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "parent" {
  name = "test-ivparent"
  type = "dummy"
  dummy {}
}
resource "iproute_link" "test" {
  name   = "test-ipvlan0"
  type   = "ipvlan"
  ipvlan {
    parent = iproute_link.parent.name
    mode   = "l2"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-ipvlan0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "ipvlan"),
				),
			},
		},
	})
}

func TestAccLinkResource_vlan(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "parent" {
  name = "test-vlparent"
  type = "dummy"
  dummy {}
}
resource "iproute_link" "test" {
  name   = "test-vlan100"
  type   = "vlan"
  vlan {
    parent  = iproute_link.parent.name
    vlan_id = 100
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-vlan100"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "vlan"),
				),
			},
		},
	})
}

func TestAccLinkResource_vti(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-vti0"
  type = "vti"
  vti {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-vti0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "vti"),
				),
			},
		},
	})
}

func TestAccLinkResource_ip6tnl(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-ip6tnl0"
  type = "ip6tnl"
  ip6tnl {
    local  = "fd00::1"
    remote = "fd00::2"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-ip6tnl0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "ip6tnl"),
				),
			},
		},
	})
}

func TestAccLinkResource_ifb(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-ifb0"
  type = "ifb"
  ifb {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-ifb0"),
					resource.TestCheckResourceAttr("iproute_link.test", "type", "ifb"),
				),
			},
		},
	})
}

func TestAccLinkResource_description(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-desc0"
  type        = "dummy"
  description = "My test interface"
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-desc0"),
					resource.TestCheckResourceAttr("iproute_link.test", "description", "My test interface"),
				),
			},
		},
	})
}

func TestAccLinkResource_txqueuelen(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name         = "test-txq0"
  type         = "dummy"
  tx_queue_len = 500
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "tx_queue_len", "500"),
				),
			},
		},
	})
}

func TestAccLinkResource_mac_address(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name        = "test-mac0"
  type        = "dummy"
  mac_address = "aa:bb:cc:dd:ee:ff"
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "mac_address", "aa:bb:cc:dd:ee:ff"),
				),
			},
		},
	})
}

func TestAccLinkResource_bridge_with_port(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "br" {
  name = "test-brp0"
  type = "bridge"
  bridge {}
}
resource "iproute_link" "port" {
  name   = "test-brport0"
  type   = "dummy"
  master = iproute_link.br.name
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.port", "master", "test-brp0"),
				),
			},
		},
	})
}

func TestAccLinkResource_gretap(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-gretap0"
  type = "gretap"
  gre {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-gretap0"),
				),
			},
		},
	})
}
