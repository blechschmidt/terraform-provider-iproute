# List all nexthop objects
data "iproute_nexthop" "all" {}

output "nexthops" {
  value = data.iproute_nexthop.all.nexthops
}
