# Route from a source subnet to a custom table
resource "iproute_rule" "from_subnet" {
  priority = 1000
  src      = "10.0.0.0/24"
  table    = 100
}

# Route to a destination subnet
resource "iproute_rule" "to_subnet" {
  priority = 1001
  dst      = "192.168.0.0/16"
  table    = 200
}

# Firewall mark based routing
resource "iproute_rule" "fwmark" {
  priority = 2000
  fwmark   = 42
  fwmask   = 0xff
  table    = 300
}

# Input interface rule
resource "iproute_rule" "iif" {
  priority = 3000
  iif_name = "eth0"
  table    = 100
}

# Output interface rule
resource "iproute_rule" "oif" {
  priority = 3001
  oif_name = "eth1"
  table    = 200
}

# TOS-based routing
resource "iproute_rule" "tos" {
  priority = 4000
  tos      = 8
  table    = 100
}

# IPv6 rule
resource "iproute_rule" "ipv6" {
  priority = 5000
  family   = "inet6"
  src      = "fd00::/64"
  table    = 100
}
