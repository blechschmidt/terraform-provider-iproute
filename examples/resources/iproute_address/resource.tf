# IPv4 address
resource "iproute_address" "ipv4" {
  address = "10.0.0.1/24"
  device  = "dummy0"
}

# IPv6 address
resource "iproute_address" "ipv6" {
  address = "fd00::1/64"
  device  = "dummy0"
  scope   = "global"
}

# Address with broadcast
resource "iproute_address" "with_broadcast" {
  address   = "192.168.1.1/24"
  device    = "eth0"
  broadcast = "192.168.1.255"
}

# Address with label
resource "iproute_address" "labeled" {
  address = "10.0.0.2/24"
  device  = "eth0"
  label   = "eth0:web"
}

# Point-to-point address with peer
resource "iproute_address" "ptp" {
  address = "10.0.0.1/32"
  device  = "tun0"
  peer    = "10.0.0.2/32"
}

# Multiple addresses on one interface
resource "iproute_address" "secondary" {
  address = "10.0.0.3/24"
  device  = "dummy0"
}
