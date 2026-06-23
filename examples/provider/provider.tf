# Default namespace
provider "iproute" {}

# Named network namespace (under /run/netns)
provider "iproute" {
  alias     = "isolated"
  namespace = "my-namespace"
}

# Process network namespace by PID
provider "iproute" {
  alias     = "by_pid"
  namespace = "pid:12345"
}

# nsfs path (e.g. a bind mount or /proc/<pid>/ns/net)
provider "iproute" {
  alias     = "by_path"
  namespace = "path:/proc/12345/ns/net"
}

# A running Docker container's network namespace, resolved via the Docker API
provider "iproute" {
  alias     = "container"
  namespace = "docker:${docker_container.gateway.id}"
}
