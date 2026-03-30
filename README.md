# Terraform Provider: iproute

A Terraform provider for managing Linux networking resources via netlink. It provides declarative management of everything you would normally configure with the `ip` command from iproute2: links, addresses, routes, rules, neighbors, tunnels, namespaces, and more.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Go >= 1.24 (for building from source)
- Linux (the provider uses netlink, which is Linux-specific)
- Root or `CAP_NET_ADMIN` capability

## Installation

### From Source

```sh
git clone https://github.com/example/terraform-provider-iproute.git
cd terraform-provider-iproute
make install
```

### Development Override

Add a dev override to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/example/iproute" = "/path/to/go/bin"
  }
  direct {}
}
```

## Provider Configuration

```hcl
provider "iproute" {
  # Optional: operate in a specific network namespace.
  # If omitted, the provider uses the default namespace.
  namespace = "my-namespace"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `namespace` | String | No | Network namespace to operate in. If not set, operates in the default namespace. |

## Resources

| Resource | Description |
|----------|-------------|
| [iproute_link](#iproute_link) | Network links/interfaces (`ip link`) |
| [iproute_address](#iproute_address) | IP addresses (`ip address`) |
| [iproute_route](#iproute_route) | Routing table entries (`ip route`) |
| [iproute_rule](#iproute_rule) | Routing policy rules (`ip rule`) |
| [iproute_neighbor](#iproute_neighbor) | Neighbor/ARP entries (`ip neigh`) |
| [iproute_netns](#iproute_netns) | Network namespaces (`ip netns`) |
| [iproute_nexthop](#iproute_nexthop) | Nexthop objects (`ip nexthop`) |
| [iproute_tunnel](#iproute_tunnel) | IP tunnels (`ip tunnel`) |
| [iproute_l2tp_tunnel](#iproute_l2tp_tunnel) | L2TP tunnels (`ip l2tp`) |
| [iproute_l2tp_session](#iproute_l2tp_session) | L2TP sessions (`ip l2tp`) |
| [iproute_fou](#iproute_fou) | Foo-over-UDP (`ip fou`) |
| [iproute_xfrm_state](#iproute_xfrm_state) | XFRM/IPsec states (`ip xfrm state`) |
| [iproute_xfrm_policy](#iproute_xfrm_policy) | XFRM/IPsec policies (`ip xfrm policy`) |
| [iproute_macsec](#iproute_macsec) | MACsec devices (`ip macsec`) |
| [iproute_tuntap](#iproute_tuntap) | TUN/TAP devices (`ip tuntap`) |
| [iproute_maddress](#iproute_maddress) | Multicast addresses (`ip maddress`) |
| [iproute_token](#iproute_token) | IPv6 token identifiers (`ip token`) |
| [iproute_tcp_metrics](#iproute_tcp_metrics) | TCP metrics cache (`ip tcp_metrics`) |
| [iproute_sr](#iproute_sr) | Segment routing (`ip sr`) |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| [iproute_link](#data-iproute_link) | Read link information |
| [iproute_address](#data-iproute_address) | Read IP addresses |
| [iproute_route](#data-iproute_route) | Read routing table |
| [iproute_rule](#data-iproute_rule) | Read routing policy rules |
| [iproute_neighbor](#data-iproute_neighbor) | Read neighbor/ARP entries |
| [iproute_netns](#data-iproute_netns) | List network namespaces |
| [iproute_nexthop](#data-iproute_nexthop) | List nexthop objects |

---

## Resource Reference

### iproute_link

Manages a network link/interface. Supports many link types including dummy, bridge, veth, vlan, vxlan, bond, macvlan, ipvlan, geneve, wireguard, GRE, SIT, IPIP, VTI, IP6TNL, tuntap, and more.

#### Example

```hcl
# Dummy interface
resource "iproute_link" "dummy" {
  name = "dummy0"
  type = "dummy"
  dummy {}
}

# Bridge with STP
resource "iproute_link" "bridge" {
  name = "br0"
  type = "bridge"
  bridge {
    stp           = true
    vlan_filtering = true
  }
}

# Veth pair
resource "iproute_link" "veth" {
  name    = "veth0"
  type    = "veth"
  enabled = true
  veth {
    peer_name = "veth1"
  }
}

# WireGuard interface
resource "iproute_link" "wg" {
  name = "wg0"
  type = "wireguard"
  mtu  = 1420
  wireguard {
    private_key = "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk="
    listen_port = 51820
    peers {
      public_key  = "xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg="
      endpoint    = "192.168.1.1:51820"
      allowed_ips = ["10.0.0.0/24"]
    }
  }
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Interface name (max 15 chars). Forces replacement. |
| `type` | String | Yes | Link type. Forces replacement. |
| `enabled` | Boolean | No | Administrative state (up/down). Default: `true`. |
| `description` | String | No | Interface description (ifAlias). |
| `mac_address` | String | No | Hardware MAC address. |
| `mtu` | Int64 | No | Maximum transmission unit. |
| `tx_queue_len` | Int64 | No | Transmit queue length. |
| `master` | String | No | Master device (e.g., bridge name). |
| `id` | String | Computed | Resource identifier. |
| `admin_status` | String | Computed | Administrative status. |
| `oper_status` | String | Computed | Operational status. |
| `if_index` | Int64 | Computed | Interface index. |
| `speed` | Int64 | Computed | Interface speed in Mbps. |
| `statistics` | Object | Computed | Interface statistics (rx/tx bytes, packets, errors, dropped, multicast). |

#### Link Type Blocks

Each link type has a corresponding configuration block. Only the block matching `type` should be provided.

<details>
<summary><strong>bridge</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `stp` | Boolean | Enable STP |
| `hello_time` | Int64 | STP hello time (centiseconds) |
| `max_age` | Int64 | STP max age (centiseconds) |
| `forward_delay` | Int64 | STP forward delay (centiseconds) |
| `vlan_filtering` | Boolean | Enable VLAN filtering |
| `default_pvid` | Int64 | Default PVID |
| `ageing_time` | Int64 | MAC ageing time (centiseconds) |
</details>

<details>
<summary><strong>bond</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `mode` | String | Bonding mode (balance-rr, active-backup, balance-xor, broadcast, 802.3ad, balance-tlb, balance-alb) |
| `miimon` | Int64 | MII monitoring interval (ms) |
| `up_delay` | Int64 | Delay before enabling slave (ms) |
| `down_delay` | Int64 | Delay before disabling slave (ms) |
| `primary` | String | Primary slave interface |
| `lacp_rate` | String | LACP rate (slow, fast) |
| `xmit_hash_policy` | String | Transmit hash policy |
</details>

<details>
<summary><strong>veth</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `peer_name` | String | Peer interface name |
</details>

<details>
<summary><strong>vlan</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `vlan_id` | Int64 | VLAN ID (1-4094) |
| `parent` | String | Parent interface name |
| `protocol` | String | VLAN protocol (802.1Q, 802.1ad) |
</details>

<details>
<summary><strong>vxlan</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `vni` | Int64 | VxLAN Network Identifier |
| `group` | String | Multicast group address |
| `local` | String | Local IP address |
| `dev` | String | Physical device for endpoints |
| `port` | Int64 | UDP destination port |
| `learning` | Boolean | Enable MAC learning |
| `proxy` | Boolean | Enable ARP proxy |
| `l2miss` | Boolean | Enable LLADDR miss notifications |
| `l3miss` | Boolean | Enable IP addr miss notifications |
</details>

<details>
<summary><strong>macvlan / macvtap</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `parent` | String | Parent interface name |
| `mode` | String | Mode (private, vepa, bridge, passthru, source) |
</details>

<details>
<summary><strong>ipvlan</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `parent` | String | Parent interface name |
| `mode` | String | Mode (l2, l3, l3s) |
| `flag` | String | Flag (bridge, private, vepa) |
</details>

<details>
<summary><strong>geneve</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `vni` | Int64 | Virtual Network Identifier |
| `remote` | String | Remote IP address |
| `port` | Int64 | UDP port |
| `ttl` | Int64 | TTL |
| `tos` | Int64 | TOS |
</details>

<details>
<summary><strong>gre</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `local` | String | Local endpoint IP |
| `remote` | String | Remote endpoint IP |
| `ttl` | Int64 | TTL |
| `tos` | Int64 | TOS |
| `pmtu_disc` | Boolean | Enable PMTU discovery |
| `key` | Int64 | Tunnel key |
| `ikey` | Int64 | Input key |
| `okey` | Int64 | Output key |
</details>

<details>
<summary><strong>sit</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `local` | String | Local endpoint IP |
| `remote` | String | Remote endpoint IP |
| `ttl` | Int64 | TTL |
| `tos` | Int64 | TOS |
</details>

<details>
<summary><strong>ipip</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `local` | String | Local endpoint IP |
| `remote` | String | Remote endpoint IP |
| `ttl` | Int64 | TTL |
| `tos` | Int64 | TOS |
</details>

<details>
<summary><strong>vti</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `local` | String | Local endpoint IP |
| `remote` | String | Remote endpoint IP |
| `ikey` | Int64 | Input key |
| `okey` | Int64 | Output key |
</details>

<details>
<summary><strong>ip6tnl</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `local` | String | Local IPv6 endpoint |
| `remote` | String | Remote IPv6 endpoint |
| `ttl` | Int64 | TTL/Hop limit |
| `flow_label` | Int64 | Flow label |
| `encap_limit` | Int64 | Encap limit |
</details>

<details>
<summary><strong>wireguard</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `private_key` | String (sensitive) | Private key |
| `listen_port` | Int64 | Listen port |
| `fwmark` | Int64 | Firewall mark |
| `peers` | List (nested) | WireGuard peers |

Peers attributes:

| Name | Type | Description |
|------|------|-------------|
| `public_key` | String | Peer public key (required) |
| `preshared_key` | String (sensitive) | Preshared key |
| `endpoint` | String | Endpoint (host:port) |
| `allowed_ips` | List(String) | Allowed IP ranges |
| `persistent_keepalive` | Int64 | Keepalive interval (seconds) |
</details>

<details>
<summary><strong>tuntap</strong></summary>

| Name | Type | Description |
|------|------|-------------|
| `mode` | String | Mode (tun, tap) |
| `owner` | Int64 | Owner UID |
| `group` | Int64 | Group GID |
| `multi_queue` | Boolean | Enable multi-queue |
</details>

<details>
<summary><strong>dummy / ifb</strong></summary>

No configuration attributes. Just provide the empty block:

```hcl
resource "iproute_link" "dummy" {
  name = "dummy0"
  type = "dummy"
  dummy {}
}
```
</details>

#### Import

```sh
terraform import iproute_link.example <interface-name>
```

---

### iproute_address

Manages an IP address on a network interface.

#### Example

```hcl
resource "iproute_address" "ipv4" {
  address = "10.0.0.1/24"
  device  = iproute_link.dummy.name
}

resource "iproute_address" "ipv6" {
  address = "fd00::1/64"
  device  = iproute_link.dummy.name
  scope   = "global"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | String | Yes | IP address in CIDR notation. Forces replacement. |
| `device` | String | Yes | Interface name. Forces replacement. |
| `peer` | String | No | Peer address for point-to-point interfaces. |
| `broadcast` | String | No | Broadcast address. |
| `label` | String | No | Address label. |
| `scope` | String | No | Address scope (global, link, host, site). |
| `flags` | List(String) | No | Address flags. |
| `preferred_lifetime` | Int64 | No | Preferred lifetime in seconds. |
| `valid_lifetime` | Int64 | No | Valid lifetime in seconds. |
| `id` | String | Computed | Resource identifier. |
| `family` | String | Computed | Address family (inet or inet6). |
| `origin` | String | Computed | Address origin. |

#### Import

```sh
terraform import iproute_address.example <device>|<address-cidr>
```

---

### iproute_route

Manages a routing table entry.

#### Example

```hcl
resource "iproute_route" "default" {
  destination = "default"
  gateway     = "10.0.0.1"
  device      = "eth0"
}

resource "iproute_route" "subnet" {
  destination = "192.168.1.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.dummy.name
  metric      = 100
}

resource "iproute_route" "blackhole" {
  destination = "10.99.0.0/16"
  type        = "blackhole"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `destination` | String | Yes | Destination in CIDR notation or `default`. Forces replacement. |
| `gateway` | String | No | Gateway IP address. |
| `device` | String | No | Output device name. |
| `source` | String | No | Source address for route. |
| `metric` | Int64 | No | Route metric/priority. |
| `table` | Int64 | No | Routing table ID. |
| `scope` | String | No | Route scope. |
| `protocol` | String | No | Route protocol. |
| `type` | String | No | Route type (unicast, blackhole, unreachable, prohibit). |
| `mtu` | Int64 | No | Route MTU. |
| `advmss` | Int64 | No | Advertised MSS. |
| `nexthop_id` | Int64 | No | Nexthop ID. |
| `onlink` | Boolean | No | Treat gateway as on-link. |
| `multipath` | List(Object) | No | ECMP multipath routes (`gateway`, `device`, `weight`). |
| `id` | String | Computed | Resource identifier. |
| `family` | String | Computed | Address family. |

#### Import

```sh
terraform import iproute_route.example "<destination>|<table>"
# Example: terraform import iproute_route.example "10.1.0.0/24|main"
```

---

### iproute_rule

Manages a routing policy rule.

#### Example

```hcl
resource "iproute_rule" "from_subnet" {
  priority = 1000
  src      = "10.0.0.0/24"
  table    = 100
}

resource "iproute_rule" "fwmark" {
  priority = 2000
  fwmark   = 42
  table    = 200
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `priority` | Int64 | No | Rule priority. |
| `family` | String | No | Address family. |
| `src` | String | No | Source prefix. |
| `dst` | String | No | Destination prefix. |
| `iif_name` | String | No | Input interface. |
| `oif_name` | String | No | Output interface. |
| `table` | Int64 | No | Routing table. |
| `fwmark` | Int64 | No | Firewall mark. |
| `fwmask` | Int64 | No | Firewall mark mask. |
| `tos` | Int64 | No | TOS value. |
| `action` | String | No | Rule action. |
| `goto_priority` | Int64 | No | Goto priority. |
| `suppress_prefix_len` | Int64 | No | Suppress prefix length. |
| `suppress_if_group` | Int64 | No | Suppress if group. |
| `ip_proto` | Int64 | No | IP protocol. |
| `sport_range` | String | No | Source port range (start-end). |
| `dport_range` | String | No | Destination port range (start-end). |
| `uid_range` | String | No | UID range (start-end). |
| `invert` | Boolean | No | Invert match. |
| `id` | String | Computed | Resource identifier. |

#### Import

```sh
terraform import iproute_rule.example <priority>
```

---

### iproute_neighbor

Manages a neighbor/ARP entry.

#### Example

```hcl
resource "iproute_neighbor" "static" {
  address = "10.0.0.2"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = iproute_link.dummy.name
  state   = "permanent"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | String | Yes | IP address. Forces replacement. |
| `lladdr` | String | Yes | Link-layer address (MAC). |
| `device` | String | Yes | Interface name. Forces replacement. |
| `state` | String | No | Neighbor state (permanent, noarp, reachable, stale). |
| `proxy` | Boolean | No | Add proxy ARP entry. |
| `vni` | Int64 | No | VNI for FDB entries. |
| `flags` | List(String) | No | Neighbor flags. |
| `id` | String | Computed | Resource identifier. |
| `family` | String | Computed | Address family. |
| `is_router` | Boolean | Computed | Whether neighbor is a router. |
| `origin` | String | Computed | Neighbor origin. |

#### Import

```sh
terraform import iproute_neighbor.example <device>|<address>
```

---

### iproute_netns

Manages a network namespace.

#### Example

```hcl
resource "iproute_netns" "isolated" {
  name = "isolated"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Namespace name. Forces replacement. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_nexthop

Manages a nexthop object.

#### Example

```hcl
resource "iproute_nexthop" "gw" {
  nhid    = 1
  gateway = "10.0.0.1"
  device  = iproute_link.dummy.name
}

resource "iproute_nexthop" "blackhole" {
  nhid      = 2
  blackhole = true
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `nhid` | Int64 | Yes | Nexthop ID. |
| `gateway` | String | No | Gateway IP. |
| `device` | String | No | Output device. |
| `blackhole` | Boolean | No | Blackhole nexthop. |
| `family` | String | No | Address family. |
| `group` | List(Object) | No | Nexthop group members (`id`, `weight`). |
| `resilient` | Boolean | No | Resilient nexthop group. |
| `fdb` | Boolean | No | FDB nexthop. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_tunnel

Manages an IP tunnel. A simplified wrapper around link types GRE, SIT, and IPIP.

#### Example

```hcl
resource "iproute_tunnel" "gre" {
  name   = "gre1"
  mode   = "gre"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
  ttl    = 64
}

resource "iproute_tunnel" "ipip" {
  name   = "ipip1"
  mode   = "ipip"
  local  = "10.0.0.1"
  remote = "10.0.0.2"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Tunnel interface name. |
| `mode` | String | Yes | Tunnel mode (gre, sit, ipip). Forces replacement. |
| `local` | String | No | Local endpoint. |
| `remote` | String | No | Remote endpoint. |
| `ttl` | Int64 | No | TTL. |
| `tos` | Int64 | No | TOS. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_l2tp_tunnel

Manages an L2TP tunnel.

#### Example

```hcl
resource "iproute_l2tp_tunnel" "example" {
  tunnel_id      = 1
  peer_tunnel_id = 1
  encap_type     = "udp"
  local          = "10.0.0.1"
  remote         = "10.0.0.2"
  local_port     = 5000
  remote_port    = 5000
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `tunnel_id` | Int64 | Yes | Tunnel ID. |
| `peer_tunnel_id` | Int64 | Yes | Peer tunnel ID. |
| `encap_type` | String | Yes | Encapsulation type (udp, ip). |
| `local` | String | Yes | Local IP address. |
| `remote` | String | Yes | Remote IP address. |
| `local_port` | Int64 | No | Local UDP port. |
| `remote_port` | Int64 | No | Remote UDP port. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_l2tp_session

Manages an L2TP session.

#### Example

```hcl
resource "iproute_l2tp_session" "example" {
  tunnel_id       = iproute_l2tp_tunnel.example.tunnel_id
  session_id      = 1
  peer_session_id = 1
  name            = "l2tp-sess0"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `tunnel_id` | Int64 | Yes | Tunnel ID. |
| `session_id` | Int64 | Yes | Session ID. |
| `peer_session_id` | Int64 | Yes | Peer session ID. |
| `name` | String | No | Session interface name. |
| `cookie` | String | No | Cookie value. |
| `peer_cookie` | String | No | Peer cookie value. |
| `l2spec_type` | String | No | L2-specific header type. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_fou

Manages a Foo-over-UDP (FOU) receive port.

#### Example

```hcl
resource "iproute_fou" "direct" {
  port     = 5555
  protocol = 4  # IPIP
}

resource "iproute_fou" "gue" {
  port       = 6666
  encap_type = "gue"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `port` | Int64 | Yes | UDP port number. |
| `family` | String | No | Address family. |
| `protocol` | Int64 | No | IP protocol number. |
| `encap_type` | String | No | Encapsulation type (direct, gue). |
| `remote_port` | Int64 | No | Remote port. |
| `local` | String | No | Local address. |
| `peer` | String | No | Peer address. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_xfrm_state

Manages an XFRM/IPsec security association state.

#### Example

```hcl
resource "iproute_xfrm_state" "esp" {
  src   = "10.0.0.1"
  dst   = "10.0.0.2"
  proto = "esp"
  spi   = 5000
  mode  = "tunnel"

  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "0x0123456789abcdef0123456789abcdef01234567"
    icv_len = 128
  }
}

resource "iproute_xfrm_state" "ah" {
  src   = "10.0.0.1"
  dst   = "10.0.0.2"
  proto = "ah"
  spi   = 6000

  auth {
    name = "hmac(sha256)"
    key  = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  }
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `src` | String | Yes | Source IP. |
| `dst` | String | Yes | Destination IP. |
| `proto` | String | Yes | Protocol (esp, ah, comp). |
| `spi` | Int64 | Yes | Security Parameter Index. |
| `mode` | String | No | Mode (transport, tunnel). |
| `reqid` | Int64 | No | Request ID. |
| `replay_window` | Int64 | No | Replay window size. |
| `mark` | Int64 | No | Mark value. |
| `mark_mask` | Int64 | No | Mark mask. |
| `if_id` | Int64 | No | Interface ID. |
| `id` | String | Computed | Resource identifier. |
| `family` | String | Computed | Address family. |

**Blocks** (provide one of `auth`, `crypt`, or `aead`):

| Block | Attributes |
|-------|------------|
| `auth` | `name` (String), `key` (String, sensitive) |
| `crypt` | `name` (String), `key` (String, sensitive) |
| `aead` | `name` (String), `key` (String, sensitive), `icv_len` (Int64) |

---

### iproute_xfrm_policy

Manages an XFRM/IPsec policy.

#### Example

```hcl
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

resource "iproute_xfrm_policy" "block" {
  src    = "10.0.0.0/24"
  dst    = "10.99.0.0/24"
  dir    = "out"
  action = "block"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `src` | String | Yes | Source prefix. |
| `dst` | String | Yes | Destination prefix. |
| `dir` | String | Yes | Direction (in, out, fwd). |
| `priority` | Int64 | No | Policy priority. |
| `action` | String | No | Action (allow, block). |
| `proto` | String | No | Protocol. |
| `src_port` | Int64 | No | Source port. |
| `dst_port` | Int64 | No | Destination port. |
| `mark` | Int64 | No | Mark value. |
| `mark_mask` | Int64 | No | Mark mask. |
| `if_id` | Int64 | No | Interface ID. |
| `templates` | List(Object) | No | Policy templates (`src`, `dst`, `proto`, `mode`, `reqid`, `spi`). |
| `id` | String | Computed | Resource identifier. |
| `family` | String | Computed | Address family. |

---

### iproute_macsec

Manages a MACsec device.

#### Example

```hcl
resource "iproute_macsec" "secure" {
  parent  = "eth0"
  name    = "macsec0"
  encrypt = true
  port    = 1
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `parent` | String | Yes | Parent interface. |
| `name` | String | Yes | MACsec interface name. |
| `sci` | String | No | SCI value. |
| `port` | Int64 | No | Port number. |
| `encrypt` | Boolean | No | Enable encryption. |
| `cipher_suite` | String | No | Cipher suite. |
| `icv_len` | Int64 | No | ICV length. |
| `encoding_sa` | Int64 | No | Encoding SA. |
| `validate` | String | No | Validate mode. |
| `protect` | Boolean | No | Enable protection. |
| `replay_protect` | Boolean | No | Enable replay protection. |
| `window` | Int64 | No | Replay window size. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_tuntap

Manages a TUN/TAP device.

#### Example

```hcl
resource "iproute_tuntap" "tun" {
  name = "tun0"
  mode = "tun"
}

resource "iproute_tuntap" "tap" {
  name        = "tap0"
  mode        = "tap"
  multi_queue = true
  owner       = 1000
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Device name. |
| `mode` | String | Yes | Mode (tun, tap). |
| `owner` | Int64 | No | Owner UID. |
| `group` | Int64 | No | Group GID. |
| `multi_queue` | Boolean | No | Enable multi-queue. |
| `persist` | Boolean | No | Persistent device. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_maddress

Manages a multicast address on an interface.

#### Example

```hcl
resource "iproute_maddress" "multicast" {
  device  = iproute_link.dummy.name
  address = "33:33:00:00:00:01"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `device` | String | Yes | Interface name. |
| `address` | String | Yes | Multicast address. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_token

Manages IPv6 tokenized interface identifiers.

#### Example

```hcl
resource "iproute_token" "example" {
  device = iproute_link.veth.name
  token  = "::1"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `device` | String | Yes | Interface name (must support IPv6 NDP). |
| `token` | String | Yes | IPv6 interface token. |
| `id` | String | Computed | Resource identifier. |

---

### iproute_tcp_metrics

Manages TCP metrics cache entries. Note: TCP metrics entries are typically created by the kernel through actual TCP connections. This resource manages existing cache entries.

#### Example

```hcl
resource "iproute_tcp_metrics" "example" {
  address = "10.0.0.1"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `address` | String | Yes | Destination address. |
| `id` | String | Computed | Resource identifier. |
| `rtt` | Int64 | Computed | RTT (microseconds). |
| `rttvar` | Int64 | Computed | RTT variance. |
| `ssthresh` | Int64 | Computed | Slow start threshold. |
| `cwnd` | Int64 | Computed | Congestion window. |

---

### iproute_sr

Manages segment routing configuration.

#### Example

```hcl
resource "iproute_sr" "example" {
  device   = iproute_link.dummy.name
  segments = ["fc00::1", "fc00::2"]
  encap    = "encap"
}
```

#### Attributes

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `device` | String | No | Interface name. |
| `hmac` | String | No | HMAC key. |
| `segments` | List(String) | No | SRv6 segments. |
| `encap` | String | No | Encap mode (inline, encap). |
| `id` | String | Computed | Resource identifier. |

---

## Data Source Reference

### Data: iproute_link

Read information about a network link.

```hcl
data "iproute_link" "eth0" {
  name = "eth0"
}

output "eth0_mtu" {
  value = data.iproute_link.eth0.mtu
}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Interface name. |
| `id` | String | Computed | Resource identifier. |
| `type` | String | Computed | Link type. |
| `mtu` | Int64 | Computed | MTU. |
| `mac_address` | String | Computed | MAC address. |
| `admin_status` | String | Computed | Administrative status. |
| `oper_status` | String | Computed | Operational status. |
| `if_index` | Int64 | Computed | Interface index. |
| `tx_queue_len` | Int64 | Computed | Transmit queue length. |
| `master` | String | Computed | Master interface. |

---

### Data: iproute_address

Read IP addresses on an interface.

```hcl
data "iproute_address" "eth0" {
  device = "eth0"
  family = "inet"
}

output "addresses" {
  value = data.iproute_address.eth0.addresses
}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `device` | String | Yes | Interface name. |
| `family` | String | No | Filter by family (inet, inet6). |
| `id` | String | Computed | Resource identifier. |
| `addresses` | List(String) | Computed | Addresses in CIDR notation. |

---

### Data: iproute_route

Read routing table entries.

```hcl
data "iproute_route" "all" {
  family = "inet"
}

data "iproute_route" "custom_table" {
  table = 100
}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `family` | String | No | Address family (inet, inet6). |
| `table` | Int64 | No | Routing table ID. |
| `id` | String | Computed | Resource identifier. |
| `routes` | List(String) | Computed | List of routes. |

---

### Data: iproute_rule

Read routing policy rules.

```hcl
data "iproute_rule" "all" {
  family = "inet"
}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `family` | String | No | Address family (inet, inet6). |
| `id` | String | Computed | Resource identifier. |
| `rules` | List(String) | Computed | List of rules. |

---

### Data: iproute_neighbor

Read neighbor/ARP entries.

```hcl
data "iproute_neighbor" "eth0" {
  device = "eth0"
}

data "iproute_neighbor" "all" {}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `device` | String | No | Filter by interface. |
| `id` | String | Computed | Resource identifier. |
| `neighbors` | List(String) | Computed | List of neighbor entries. |

---

### Data: iproute_netns

List network namespaces.

```hcl
data "iproute_netns" "all" {}

output "namespaces" {
  value = data.iproute_netns.all.namespaces
}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Computed | Resource identifier. |
| `namespaces` | List(String) | Computed | List of namespace names. |

---

### Data: iproute_nexthop

List nexthop objects.

```hcl
data "iproute_nexthop" "all" {}
```

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Computed | Resource identifier. |
| `nexthops` | List(String) | Computed | List of nexthop entries. |

---

## Full Example

A complete example setting up an isolated network with a bridge, addresses, routes, and policy rules:

```hcl
provider "iproute" {}

# Create a bridge
resource "iproute_link" "br0" {
  name = "br0"
  type = "bridge"
  bridge {
    stp = true
  }
}

# Create veth pairs and attach to bridge
resource "iproute_link" "veth0" {
  name    = "veth0"
  type    = "veth"
  enabled = true
  master  = iproute_link.br0.name
  veth {
    peer_name = "veth1"
  }
}

# Assign an address to the bridge
resource "iproute_address" "br0_addr" {
  address = "10.0.0.1/24"
  device  = iproute_link.br0.name
}

# Add a static neighbor
resource "iproute_neighbor" "server" {
  address = "10.0.0.10"
  lladdr  = "aa:bb:cc:dd:ee:ff"
  device  = iproute_link.br0.name
  state   = "permanent"
  depends_on = [iproute_address.br0_addr]
}

# Add a route
resource "iproute_route" "backend" {
  destination = "192.168.1.0/24"
  gateway     = "10.0.0.254"
  device      = iproute_link.br0.name
  depends_on  = [iproute_address.br0_addr]
}

# Add a policy rule
resource "iproute_rule" "custom" {
  priority = 1000
  src      = "10.0.0.0/24"
  table    = 100
}
```

## Building and Testing

```sh
# Build
make build

# Run acceptance tests (requires root)
sudo TF_ACC=1 go test ./internal/provider/ -v -timeout 30m

# Run unit tests
go test ./internal/netlink/ ./internal/validators/ -v
```

## License

See [LICENSE](LICENSE) for details.
