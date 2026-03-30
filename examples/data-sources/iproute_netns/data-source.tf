# List all network namespaces
data "iproute_netns" "all" {}

output "namespaces" {
  value = data.iproute_netns.all.namespaces
}
