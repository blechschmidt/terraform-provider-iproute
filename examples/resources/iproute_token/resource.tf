# IPv6 token on a veth interface
# Note: the interface must support IPv6 NDP (dummy interfaces do not)
resource "iproute_token" "example" {
  device = "veth0"
  token  = "::1"
}
