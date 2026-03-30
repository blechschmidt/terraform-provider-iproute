# L2TP session
resource "iproute_l2tp_session" "example" {
  tunnel_id       = iproute_l2tp_tunnel.udp.tunnel_id
  session_id      = 1
  peer_session_id = 1
  name            = "l2tp-sess0"
}
