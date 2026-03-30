# TUN device
resource "iproute_tuntap" "tun" {
  name = "tun0"
  mode = "tun"
}

# TAP device with multi-queue
resource "iproute_tuntap" "tap" {
  name        = "tap0"
  mode        = "tap"
  multi_queue = true
  owner       = 1000
  group       = 1000
}
