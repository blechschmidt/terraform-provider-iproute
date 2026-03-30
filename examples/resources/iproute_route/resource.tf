# Basic route via gateway
resource "iproute_route" "subnet" {
  destination = "192.168.1.0/24"
  gateway     = "10.0.0.254"
  device      = "eth0"
}

# Default route
resource "iproute_route" "default" {
  destination = "default"
  gateway     = "10.0.0.1"
  device      = "eth0"
}

# Route with metric
resource "iproute_route" "with_metric" {
  destination = "10.1.0.0/16"
  gateway     = "10.0.0.1"
  device      = "eth0"
  metric      = 100
}

# Route in a custom table
resource "iproute_route" "custom_table" {
  destination = "172.16.0.0/12"
  gateway     = "10.0.0.1"
  device      = "eth0"
  table       = 100
}

# Blackhole route
resource "iproute_route" "blackhole" {
  destination = "10.99.0.0/16"
  type        = "blackhole"
}

# Unreachable route
resource "iproute_route" "unreachable" {
  destination = "10.98.0.0/16"
  type        = "unreachable"
}

# Prohibit route
resource "iproute_route" "prohibit" {
  destination = "10.97.0.0/16"
  type        = "prohibit"
}

# Route with MTU
resource "iproute_route" "with_mtu" {
  destination = "10.2.0.0/24"
  gateway     = "10.0.0.1"
  device      = "eth0"
  mtu         = 1400
}

# Route with source address
resource "iproute_route" "with_source" {
  destination = "10.3.0.0/24"
  gateway     = "10.0.0.1"
  device      = "eth0"
  source      = "10.0.0.5"
}

# IPv6 route
resource "iproute_route" "ipv6" {
  destination = "fd01::/64"
  gateway     = "fd00::1"
  device      = "eth0"
}
