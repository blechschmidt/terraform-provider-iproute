# Read link information
data "iproute_link" "eth0" {
  name = "eth0"
}

output "eth0_mtu" {
  value = data.iproute_link.eth0.mtu
}

output "eth0_mac" {
  value = data.iproute_link.eth0.mac_address
}

output "eth0_status" {
  value = data.iproute_link.eth0.oper_status
}

# Read loopback
data "iproute_link" "lo" {
  name = "lo"
}
