# ESP state with AEAD (tunnel mode)
resource "iproute_xfrm_state" "esp_aead" {
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

# ESP state with separate encryption
resource "iproute_xfrm_state" "esp_crypt" {
  src   = "10.0.0.1"
  dst   = "10.0.0.2"
  proto = "esp"
  spi   = 5001
  mode  = "transport"

  crypt {
    name = "cbc(aes)"
    key  = "0x0123456789abcdef0123456789abcdef"
  }

  auth {
    name = "hmac(sha256)"
    key  = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  }
}

# AH state
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

# IPv6 XFRM state
resource "iproute_xfrm_state" "ipv6" {
  src   = "fd00::1"
  dst   = "fd00::2"
  proto = "esp"
  spi   = 7000
  mode  = "tunnel"

  aead {
    name    = "rfc4106(gcm(aes))"
    key     = "0x0123456789abcdef0123456789abcdef01234567"
    icv_len = 128
  }
}
