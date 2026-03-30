# Create a network namespace
resource "iproute_netns" "isolated" {
  name = "isolated"
}

# Multiple namespaces
resource "iproute_netns" "app" {
  name = "app-ns"
}

resource "iproute_netns" "db" {
  name = "db-ns"
}
