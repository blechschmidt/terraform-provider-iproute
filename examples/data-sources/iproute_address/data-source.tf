# Read all IPv4 addresses on eth0
data "iproute_address" "eth0_ipv4" {
  device = "eth0"
  family = "inet"
}

# Read all addresses (IPv4 + IPv6) on eth0
data "iproute_address" "eth0_all" {
  device = "eth0"
}

output "eth0_addresses" {
  value = data.iproute_address.eth0_all.addresses
}
