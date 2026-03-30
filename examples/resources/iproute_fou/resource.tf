# FOU direct mode (Foo-over-UDP)
resource "iproute_fou" "direct" {
  port     = 5555
  protocol = 4 # IPIP
}

# GUE mode (Generic UDP Encapsulation)
resource "iproute_fou" "gue" {
  port       = 6666
  encap_type = "gue"
}
