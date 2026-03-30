# Read all IPv4 routes
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

output "routes" {
  value = data.iproute_route.ipv4.routes
}
