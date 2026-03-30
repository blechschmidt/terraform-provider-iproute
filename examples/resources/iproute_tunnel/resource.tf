# GRE tunnel
resource "iproute_tunnel" "gre" {
  name   = "gre1"
  mode   = "gre"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
  ttl    = 64
}

# IPIP tunnel
resource "iproute_tunnel" "ipip" {
  name   = "ipip1"
  mode   = "ipip"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}

# SIT tunnel (IPv6-in-IPv4)
resource "iproute_tunnel" "sit" {
  name   = "sit1"
  mode   = "sit"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
