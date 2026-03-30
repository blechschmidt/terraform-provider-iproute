# Basic MACsec device
resource "iproute_macsec" "basic" {
  parent  = "eth0"
  name    = "macsec0"
  encrypt = true
}

# MACsec with port
resource "iproute_macsec" "with_port" {
  parent  = "eth0"
  name    = "macsec1"
  encrypt = true
  port    = 1
}
