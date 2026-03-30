# TCP metrics cache entry
# Note: TCP metrics entries are typically created by the kernel via
# actual TCP connections. This resource manages existing cache entries.
resource "iproute_tcp_metrics" "example" {
  address = "10.0.0.1"
}
