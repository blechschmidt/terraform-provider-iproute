# Segment routing with SRv6 segments
resource "iproute_sr" "example" {
  device   = "eth0"
  segments = ["fc00::1", "fc00::2"]
  encap    = "encap"
}

# Inline encapsulation
resource "iproute_sr" "inline" {
  device   = "eth0"
  segments = ["fc00::10"]
  encap    = "inline"
}
