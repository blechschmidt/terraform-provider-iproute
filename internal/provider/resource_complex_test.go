package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Tests involving multiple resources working together

func TestAccComplex_linkAddressRoute(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-lar0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "test" {
  destination = "10.1.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-lar0"),
					resource.TestCheckResourceAttr("iproute_address.test", "address", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("iproute_route.test", "destination", "10.1.0.0/24"),
				),
			},
		},
	})
}

func TestAccComplex_linkAddressNeighbor(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-lan0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_neighbor" "test" {
  address    = "10.0.0.2"
  lladdr     = "aa:bb:cc:dd:ee:ff"
  device     = iproute_link.test.name
  state      = "permanent"
  depends_on = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "name", "test-lan0"),
					resource.TestCheckResourceAttr("iproute_neighbor.test", "address", "10.0.0.2"),
				),
			},
		},
	})
}

func TestAccComplex_bridgeWithAddress(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "br" {
  name = "test-brwa0"
  type = "bridge"
  bridge {}
}
resource "iproute_address" "br" {
  address = "10.0.0.1/24"
  device  = iproute_link.br.name
}
resource "iproute_link" "port" {
  name   = "test-brwa1"
  type   = "dummy"
  master = iproute_link.br.name
  dummy {}
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_address.br", "address", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("iproute_link.port", "master", "test-brwa0"),
				),
			},
		},
	})
}

func TestAccComplex_vethWithAddresses(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-vethwa0"
  type = "veth"
  veth {
    peer_name = "test-vethwa1"
  }
}
resource "iproute_address" "a" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_link.test", "type", "veth"),
					resource.TestCheckResourceAttr("iproute_address.a", "address", "10.0.0.1/24"),
				),
			},
		},
	})
}

func TestAccComplex_ruleWithRoute(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-rr0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_rule" "test" {
  priority = 7000
  src      = "10.0.0.0/24"
  table    = 700
}
resource "iproute_route" "test" {
  destination = "10.2.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  table       = 700
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.test", "table", "700"),
					resource.TestCheckResourceAttr("iproute_route.test", "table", "700"),
				),
			},
		},
	})
}

func TestAccComplex_xfrmStatePolicy(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_xfrm_state" "test" {
  src   = "10.0.0.1"
  dst   = "10.0.0.2"
  proto = "esp"
  spi   = 5000
  mode  = "tunnel"
  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "0123456789abcdef0123456789abcdef01234567"
    icv_len = 128
  }
}
resource "iproute_xfrm_policy" "test" {
  src = "10.0.0.0/24"
  dst = "10.0.5.0/24"
  dir = "out"
  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_xfrm_state.test", "spi", "5000"),
					resource.TestCheckResourceAttr("iproute_xfrm_policy.test", "dir", "out"),
				),
			},
		},
	})
}

func TestAccComplex_multipleRoutes(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-mrt0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = iproute_link.test.name
}
resource "iproute_route" "r1" {
  destination = "10.10.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
resource "iproute_route" "r2" {
  destination = "10.11.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
resource "iproute_route" "r3" {
  destination = "10.12.0.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.test.name
  depends_on  = [iproute_address.test]
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_route.r1", "destination", "10.10.0.0/24"),
					resource.TestCheckResourceAttr("iproute_route.r2", "destination", "10.11.0.0/24"),
					resource.TestCheckResourceAttr("iproute_route.r3", "destination", "10.12.0.0/24"),
				),
			},
		},
	})
}

func TestAccComplex_multipleRules(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_rule" "r1" {
  priority = 8001
  src      = "10.1.0.0/24"
  table    = 801
}
resource "iproute_rule" "r2" {
  priority = 8002
  src      = "10.2.0.0/24"
  table    = 802
}
resource "iproute_rule" "r3" {
  priority = 8003
  dst      = "10.3.0.0/24"
  table    = 803
}
`, ns),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("iproute_rule.r1", "table", "801"),
					resource.TestCheckResourceAttr("iproute_rule.r2", "table", "802"),
					resource.TestCheckResourceAttr("iproute_rule.r3", "table", "803"),
				),
			},
		},
	})
}
