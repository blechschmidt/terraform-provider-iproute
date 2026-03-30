# Nexthop with gateway
resource "iproute_nexthop" "gw" {
  nhid    = 1
  gateway = "10.0.0.1"
  device  = "eth0"
}

# Nexthop with device only
resource "iproute_nexthop" "dev" {
  nhid   = 2
  device = "eth0"
}

# Blackhole nexthop
resource "iproute_nexthop" "blackhole" {
  nhid      = 3
  blackhole = true
}

# IPv6 nexthop
resource "iproute_nexthop" "ipv6" {
  nhid    = 4
  gateway = "fd00::1"
  device  = "eth0"
  family  = "inet6"
}
