package provider_test

import (
	"fmt"
	"testing"

	"github.com/example/terraform-provider-iproute/internal/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLinkResource_invalidType(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-invalid0"
  type = "nonexistent"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccAddressResource_invalidCIDR(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-invaddr0"
  type = "dummy"
  dummy {}
}
resource "iproute_address" "test" {
  address = "not-a-cidr"
  device  = iproute_link.test.name
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccRouteResource_invalidGateway(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-invgw0"
  type = "dummy"
  dummy {}
}
resource "iproute_route" "test" {
  destination = "10.0.0.0/24"
  gateway     = "not-a-gateway"
  device      = iproute_link.test.name
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccNeighborResource_invalidMAC(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_link" "test" {
  name = "test-invmac0"
  type = "dummy"
  dummy {}
}
resource "iproute_neighbor" "test" {
  address = "10.0.0.2"
  lladdr  = "not-a-mac"
  device  = iproute_link.test.name
  state   = "permanent"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccAddressResource_missingDevice(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_address" "test" {
  address = "10.0.0.1/24"
  device  = "nonexistent-iface-99999"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccRouteResource_missingDevice(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_route" "test" {
  destination = "10.0.0.0/24"
  device      = "nonexistent-iface-99999"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccNeighborResource_missingDevice(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_neighbor" "test" {
  address = "10.0.0.2"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = "nonexistent-iface-99999"
  state   = "permanent"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccXfrmStateResource_invalidKey(t *testing.T) {
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
  spi   = 9999
  mode  = "tunnel"
  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "not-a-hex-key"
    icv_len = 128
  }
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccTunnelResource_invalidMode(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_tunnel" "test" {
  name   = "test-tuninv0"
  mode   = "invalid-mode"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}

func TestAccFouResource_invalidPort(t *testing.T) {
	ns := testutils.CreateTestNamespace(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "iproute" { namespace = %q }
resource "iproute_fou" "test" {
  port       = 0
  encap_type = "direct"
}
`, ns),
				ExpectError: testutils.ExpectErrorRegex(".*"),
			},
		},
	})
}
