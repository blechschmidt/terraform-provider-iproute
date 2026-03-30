# Static ARP entry
resource "iproute_neighbor" "static" {
  address = "10.0.0.2"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = "eth0"
  state   = "permanent"
}

# IPv6 neighbor
resource "iproute_neighbor" "ipv6" {
  address = "fd00::2"
  lladdr  = "aa:bb:cc:dd:ee:01"
  device  = "eth0"
  state   = "permanent"
}

# Proxy ARP entry
resource "iproute_neighbor" "proxy" {
  address = "10.0.0.100"
  lladdr  = "00:00:00:00:00:00"
  device  = "eth0"
  proxy   = true
}
