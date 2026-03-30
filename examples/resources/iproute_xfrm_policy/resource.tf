# Outbound policy with template
resource "iproute_xfrm_policy" "out" {
  src = "10.0.0.0/24"
  dst = "10.0.5.0/24"
  dir = "out"

  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}

# Inbound policy
resource "iproute_xfrm_policy" "in" {
  src = "10.0.5.0/24"
  dst = "10.0.0.0/24"
  dir = "in"

  templates = [{
    src   = "10.0.0.2"
    dst   = "10.0.0.1"
    proto = "esp"
    mode  = "tunnel"
  }]
}

# Forward policy
resource "iproute_xfrm_policy" "fwd" {
  src = "10.0.5.0/24"
  dst = "10.0.0.0/24"
  dir = "fwd"

  templates = [{
    src   = "10.0.0.2"
    dst   = "10.0.0.1"
    proto = "esp"
    mode  = "tunnel"
  }]
}

# Block policy
resource "iproute_xfrm_policy" "block" {
  src    = "10.0.0.0/24"
  dst    = "10.99.0.0/24"
  dir    = "out"
  action = "block"
}

# Policy with mark
resource "iproute_xfrm_policy" "marked" {
  src       = "10.0.0.0/24"
  dst       = "10.0.10.0/24"
  dir       = "out"
  mark      = 42
  mark_mask = 0xff

  templates = [{
    src   = "10.0.0.1"
    dst   = "10.0.0.2"
    proto = "esp"
    mode  = "tunnel"
  }]
}
