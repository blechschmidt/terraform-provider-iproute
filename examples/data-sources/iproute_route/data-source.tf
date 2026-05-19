# Read all IPv4 routes (with full netlink attributes)
data "iproute_route" "ipv4" {
  family = "inet"
}

# Read IPv6 routes
data "iproute_route" "ipv6" {
  family = "inet6"
}

# Read routes from a specific table
data "iproute_route" "custom_table" {
  table = 100
}

# Look up the route used to reach a specific destination
# (equivalent to: ip route get 8.8.8.8)
data "iproute_route" "to_dns" {
  get = "8.8.8.8"
}

output "default_gateway_for_dns" {
  value = data.iproute_route.to_dns.routes[0].gateway
}

output "all_destinations" {
  value = [for r in data.iproute_route.ipv4.routes : r.destination]
}
