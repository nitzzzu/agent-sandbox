# E2B Sandbox Networking Reference

## Public URLs

Every sandbox has a public URL for accessing services running inside it. The port is the leftmost part of the hostname.

```typescript
const host = sandbox.getHost(3000)
// → "3000-i62mff4ahtrdfdkyn2esc.e2b.app"
const url = `https://${host}`
```

```python
host = sandbox.get_host(3000)
url = f'https://{host}'
```

## Internet Access Control

By default, sandboxes have full internet access. Disable with `allowInternetAccess: false` / `allow_internet_access=False`.

```typescript
// No internet
const sandbox = await Sandbox.create({ allowInternetAccess: false })
```

```python
sandbox = Sandbox.create(allow_internet_access=False)
```

Setting `allowInternetAccess` to `false` is equivalent to `network.denyOut: ['0.0.0.0/0']`.

## Fine-Grained Network Control

Use `network` config for allow/deny lists with IPs, CIDRs, or domains.

```typescript
import { Sandbox, ALL_TRAFFIC } from '@e2b/code-interpreter'

// Deny all, allow specific IPs
const sandbox = await Sandbox.create({
  network: {
    denyOut: [ALL_TRAFFIC],
    allowOut: ['1.1.1.1', '8.8.8.0/24'],
  },
})
```

```python
from e2b_code_interpreter import Sandbox, ALL_TRAFFIC

sandbox = Sandbox.create(
    network={
        "deny_out": [ALL_TRAFFIC],
        "allow_out": ["1.1.1.1", "8.8.8.0/24"],
    },
)
```

### Domain-Based Filtering

Allow traffic to specific domains. Must use `ALL_TRAFFIC` in deny list. Wildcards supported.

```typescript
const sandbox = await Sandbox.create({
  network: {
    allowOut: ['api.example.com', '*.github.com'],
    denyOut: [ALL_TRAFFIC],
  },
})
```

```python
sandbox = Sandbox.create(
    network={
        "allow_out": ["api.example.com", "*.github.com"],
        "deny_out": [ALL_TRAFFIC],
    },
)
```

**Limitations:** Domain filtering works for HTTP (port 80, Host header) and TLS (port 443, SNI). Other ports use CIDR only. UDP/QUIC/HTTP3 not supported for domain filtering. DNS nameserver `8.8.8.8` is auto-allowed when domains are used.

### Priority Rules

**Allow rules always take precedence** over deny rules. If an IP is in both lists, it is allowed.

### ALL_TRAFFIC Constant

`ALL_TRAFFIC` = `'0.0.0.0/0'` — matches all IP addresses.

## Restricting Public Access (Traffic Access Token)

By default sandbox URLs are publicly accessible. Restrict with `allowPublicTraffic: false`.

```typescript
const sandbox = await Sandbox.create({
  network: { allowPublicTraffic: false },
})

// Get the access token
console.log(sandbox.trafficAccessToken)

// Unauthenticated request → 403
const response1 = await fetch(url)

// Authenticated request → 200
const response2 = await fetch(url, {
  headers: { 'e2b-traffic-access-token': sandbox.trafficAccessToken },
})
```

```python
sandbox = Sandbox.create(
    network={"allow_public_traffic": False},
)

print(sandbox.traffic_access_token)

# Authenticated request
response = requests.get(url, headers={
    'e2b-traffic-access-token': sandbox.traffic_access_token,
})
```

## Host Masking

Customize the `Host` header sent to services inside the sandbox.

```typescript
const sandbox = await Sandbox.create({
  network: { maskRequestHost: 'localhost:${PORT}' },
})
// Requests will have Host: localhost:8080 (for port 8080)
```

```python
sandbox = Sandbox.create(
    network={"mask_request_host": "localhost:${PORT}"},
)
```

The `${PORT}` variable is replaced with the actual port number.
