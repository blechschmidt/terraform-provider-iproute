# Dummy interface
resource "iproute_link" "dummy" {
  name = "dummy0"
  type = "dummy"
  dummy {}
}

# Bridge with STP and VLAN filtering
resource "iproute_link" "bridge" {
  name = "br0"
  type = "bridge"
  bridge {
    stp            = true
    vlan_filtering = true
    hello_time     = 200
  }
}

# Veth pair
resource "iproute_link" "veth" {
  name    = "veth0"
  type    = "veth"
  enabled = true
  mtu     = 9000
  veth {
    peer_name = "veth1"
  }
}

# VLAN interface
resource "iproute_link" "vlan" {
  name = "vlan100"
  type = "vlan"
  vlan {
    parent  = iproute_link.dummy.name
    vlan_id = 100
  }
}

# Bond with LACP
resource "iproute_link" "bond" {
  name = "bond0"
  type = "bond"
  bond {
    mode      = "802.3ad"
    miimon    = 100
    lacp_rate = "fast"
  }
}

# VxLAN overlay
resource "iproute_link" "vxlan" {
  name = "vxlan10"
  type = "vxlan"
  vxlan {
    vni   = 10
    group = "239.1.1.1"
    port  = 4789
    dev   = "eth0"
  }
}

# Macvlan
resource "iproute_link" "macvlan" {
  name = "macvlan0"
  type = "macvlan"
  macvlan {
    parent = iproute_link.dummy.name
    mode   = "bridge"
  }
}

# IPvlan
resource "iproute_link" "ipvlan" {
  name = "ipvlan0"
  type = "ipvlan"
  ipvlan {
    parent = iproute_link.dummy.name
    mode   = "l3"
  }
}

# Geneve tunnel
resource "iproute_link" "geneve" {
  name = "geneve0"
  type = "geneve"
  geneve {
    vni    = 100
    remote = "10.0.0.2"
    port   = 6081
  }
}

# GRE tunnel
resource "iproute_link" "gre" {
  name = "gre1"
  type = "gre"
  gre {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
    ttl    = 64
  }
}

# SIT tunnel (IPv6-in-IPv4)
resource "iproute_link" "sit" {
  name = "sit1"
  type = "sit"
  sit {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
  }
}

# IPIP tunnel
resource "iproute_link" "ipip" {
  name = "ipip1"
  type = "ipip"
  ipip {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
    ttl    = 64
  }
}

# VTI tunnel
resource "iproute_link" "vti" {
  name = "vti1"
  type = "vti"
  vti {
    local  = "10.0.0.1"
    remote = "10.0.0.2"
    ikey   = 100
    okey   = 100
  }
}

# IPv6 tunnel
resource "iproute_link" "ip6tnl" {
  name = "ip6tnl1"
  type = "ip6tnl"
  ip6tnl {
    local  = "fd00::1"
    remote = "fd00::2"
  }
}

# WireGuard
resource "iproute_link" "wireguard" {
  name = "wg0"
  type = "wireguard"
  mtu  = 1420
  wireguard {
    private_key = "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk="
    listen_port = 51820
    peers {
      public_key  = "xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg="
      endpoint    = "192.168.1.1:51820"
      allowed_ips = ["10.0.0.0/24", "10.0.1.0/24"]
      persistent_keepalive = 25
    }
  }
}

# TUN/TAP via link resource
resource "iproute_link" "tuntap" {
  name = "tap0"
  type = "tuntap"
  tuntap {
    mode        = "tap"
    multi_queue = true
  }
}

# IFB interface
resource "iproute_link" "ifb" {
  name = "ifb0"
  type = "ifb"
  ifb {}
}

# Bridge with port attached
resource "iproute_link" "br_port" {
  name    = "dummy1"
  type    = "dummy"
  master  = iproute_link.bridge.name
  enabled = true
  dummy {}
}
