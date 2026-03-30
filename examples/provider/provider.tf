# Default namespace
provider "iproute" {}

# Specific network namespace
provider "iproute" {
  alias     = "isolated"
  namespace = "my-namespace"
}
