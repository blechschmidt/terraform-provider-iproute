# Read all IPv4 routing policy rules
data "iproute_rule" "ipv4" {
  family = "inet"
}

# Read IPv6 rules
data "iproute_rule" "ipv6" {
  family = "inet6"
}

output "rules" {
  value = data.iproute_rule.ipv4.rules
}
