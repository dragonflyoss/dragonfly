# gRPC Token Authentication (JWT)

`pkg/rpc/auth` adds HS256 JWT authentication to the gRPC calls exchanged
between dragonfly components (manager, scheduler, dfdaemon). It is disabled by
default so existing deployments keep working; each component opts in through
its configuration.

## What it provides

- `TokenManager` — issues and verifies HS256 tokens
  (`IssueToken(subject, ttl)`, `ParseToken(raw)`).
- `TokenProvider` — issues a token once, caches it, and refreshes it at 80% of
  its lifetime, so high-frequency callers do not sign on every RPC.
- `UnaryServerInterceptor` / `StreamServerInterceptor` — read the
  `authorization` metadata value (`Bearer <jwt>`), verify the signature,
  algorithm, issuer, audience (optional) and expiration, and reject invalid
  requests with `codes.Unauthenticated`.
- `NewUnaryClientInterceptor` / `NewStreamClientInterceptor` — attach the
  provider's token to outgoing calls.
- `SubjectFromContext` / `ClaimsFromContext` — expose the verified identity to
  service handlers.

## Configuration

The signing key must come from component configuration or a mounted secret,
never from source code. The YAML paths differ slightly between components:

**Manager** — under `server.grpc.auth`:

```yaml
server:
  grpc:
    auth:
      enabled: true
      secret: <shared HS256 key>
      issuer: dragonfly
      tokenTTL: 24h
```

**Scheduler** — under `server.auth`:

```yaml
server:
  auth:
    enabled: true
    secret: <shared HS256 key>
    issuer: dragonfly
    tokenTTL: 24h
```

Both components perform the same validation when `enabled` is `true`: `secret`
and `issuer` must be non-empty and `tokenTTL` must be positive. Defaults:
issuer `dragonfly`, TTL `24h`.

> All components that talk to each other must share the same `secret` and
> `issuer`. The secret should be provisioned like any other credential — e.g.
> a Kubernetes Secret mounted into the containers — and rotated as a whole
> (see "Key rotation" below).

`skippedMethods` (optional, on both) lists gRPC methods that bypass
authentication: either an exact full method name
(`/grpc.health.v1.Health/Check`) or a service prefix ending in `/`
(`/grpc.health.v1.Health/`). The gRPC health checking service is always
exempt, so Kubernetes liveness/readiness probes are unaffected. Keep this list
empty unless you have a specific reason; a broad prefix such as `/` disables
authentication for the entire server.

The `Audience` field of `auth.Config` exists at the Go API level but is not
exposed in component configuration; when unset, the audience claim is neither
issued nor checked.

## How the components use it

When enabled, `manager` and `scheduler` install the server interceptors in
front of their existing interceptor chains (rate limiting, metrics, logging,
recovery). The scheduler additionally attaches a client interceptor to its
connection to the manager: `scheduler/scheduler.go` builds a `TokenProvider`
with the same secret and issuer, and every outgoing call carries the cached
token. Tokens are re-issued after 80% of the TTL has elapsed; the remaining
20% is a safety margin for clock skew between components. Concurrent calls
hitting the refresh point share a single new signature.

Handlers can read the caller identity:

```go
subject, ok := auth.SubjectFromContext(ctx)
```

`dfdaemon` (Rust client) authentication is configured separately in the
[client repository](https://github.com/dragonflyoss/client).

## Wiring manually (library usage)

```go
tm, err := auth.NewTokenManager(auth.Config{
    Key:    secret,   // from configuration or a mounted secret
    Issuer: "dragonfly",
})
if err != nil {
    return err
}

// Server
grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
    auth.UnaryServerInterceptor(tm, "/grpc.health.v1.Health/"),
    // ... existing interceptors
))

// Client
provider := auth.NewTokenProvider(tm, "scheduler", time.Minute)
conn, err := grpc.NewClient(addr,
    grpc.WithTransportCredentials(credentials),
    grpc.WithUnaryInterceptor(auth.NewUnaryClientInterceptor(provider)),
    grpc.WithStreamInterceptor(auth.NewStreamClientInterceptor(provider)),
)
```

A token can also be attached manually with
`auth.AppendTokenToOutgoingContext(ctx, token)`.

## Behavior

| Case                                               | Result                      |
| -------------------------------------------------- | --------------------------- |
| No `authorization` metadata                        | `codes.Unauthenticated`     |
| Non-bearer scheme / empty token                    | `codes.Unauthenticated`     |
| Expired or not-yet-valid token                     | `codes.Unauthenticated`     |
| Signature signed with another key                  | `codes.Unauthenticated`     |
| Algorithm other than HS256 (including `none`)      | `codes.Unauthenticated`     |
| Issuer (or configured audience) mismatch           | `codes.Unauthenticated`     |
| Token without `exp`                                | `codes.Unauthenticated`     |
| Valid token                                        | Request reaches the handler |

## Key rotation

Rotation currently means changing the shared `secret` and restarting every
component; there is no dual-key overlap, so plan a brief maintenance window
(no new tokens are accepted during the restart). Support for a
`previousKey` overlap window is tracked as a follow-up.

## Tests

Unit tests start a real gRPC server on `bufconn` with the interceptors
installed and use the standard gRPC health service as the target:

```bash
go test ./pkg/rpc/auth/...
```

Covered cases: missing token, expired token, wrong signing key, wrong issuer,
`none` algorithm, malformed header, valid token (unary and stream), health
service exemption, claim propagation to handlers, token caching and
concurrent-refresh sharing, and expiry recovery.
