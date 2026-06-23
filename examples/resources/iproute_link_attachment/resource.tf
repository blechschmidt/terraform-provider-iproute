# Bridge an existing physical NIC into a bridge in the provider's namespace.
resource "iproute_link_attachment" "uplink" {
  name   = "eth1"
  master = iproute_link.br0.name
  up     = true
}

# Move a container's veth into the container's network namespace and enslave
# it to a bridge that already exists there. The netns argument accepts the
# same forms as the provider's namespace argument (name, name:, pid:, path:,
# docker:).
resource "iproute_link_attachment" "container_port" {
  name   = "veth-app"
  netns  = "docker:${docker_container.app.id}"
  master = "br0"
  up     = true
}
