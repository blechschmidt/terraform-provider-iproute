# Read all neighbors
data "iproute_neighbor" "all" {}

# Read neighbors on a specific interface
data "iproute_neighbor" "eth0" {
  device = "eth0"
}

output "neighbors" {
  value = data.iproute_neighbor.all.neighbors
}
