# L2TP tunnel over UDP
resource "iproute_l2tp_tunnel" "udp" {
  tunnel_id      = 1
  peer_tunnel_id = 1
  encap_type     = "udp"
  local          = "10.0.0.1"
  remote         = "10.0.0.2"
  local_port     = 5000
  remote_port    = 5000
}

# L2TP tunnel over IP
resource "iproute_l2tp_tunnel" "ip" {
  tunnel_id      = 2
  peer_tunnel_id = 2
  encap_type     = "ip"
  local          = "10.0.0.1"
  remote         = "10.0.0.2"
}
