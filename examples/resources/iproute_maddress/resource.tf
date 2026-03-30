# Multicast address
resource "iproute_maddress" "example" {
  device  = "eth0"
  address = "33:33:00:00:00:01"
}
