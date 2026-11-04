# gRPC Token Authentication (JWT)

`pkg/rpc/auth` adds HS256 JWT authentication to the gRPC calls exchanged
between dragonfly components (manager, scheduler, dfdaemon). It is a standalone
module: the existing `pkg/rpc/manager/server` and `pkg/rpc/scheduler/server`
constructors keep working unchanged until the interceptors are wired in.

## What it provides

- `TokenManager` — issues and verifies HS256 tokens
  (`IssueToken(subject, ttl)`, `ParseToken(raw)`).
- `UnaryServerInterceptor` / `StreamServerInterceptor` — read the `authorization`
  metadata value (`Bearer <jwt>`), verify the signature, algorithm, issuer,
  audience (optional) and expiration, and reject invalid requests with
  `codes.Unauthenticated`.
- `NewUnaryClientInterceptor` / `NewStreamClientInterceptor` — issue and attach
  a token for outgoing calls.
- `SubjectFromContext` / `ClaimsFromContext` — expose the verified identity to
  service handlers.

## Configuration

The signing key must come from component configuration or a mounted secret, not
from source code:

```go
tm, err := auth.NewTokenManager(auth.Config{
    Key:      cfg.Auth.GRPCJWT.Key,      // e.g. DRAGONFLY_GRPC_JWT_KEY
    Issuer:   cfg.Auth.GRPCJWT.Issuer,   // e.g. "dragonfly"
    Audience: cfg.Auth.GRPCJWT.Audience, // optional
})
if err != nil {
    return err
}
```

All components that talk to each other must share the same `Key` and `Issuer`.

## Wiring the server

Add the interceptors to the existing interceptor chains in
`pkg/rpc/manager/server/server.go` and `pkg/rpc/scheduler/server/server.go`:

```go
grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
    auth.UnaryServerInterceptor(tm, "/grpc.health.v1.Health/"),
    // ... existing interceptors
)),
grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
    auth.StreamServerInterceptor(tm, "/grpc.health.v1.Health/"),
    // ... existing interceptors
)),
```

Exempt the health checking service so that liveness and readiness probes keep
working. Any method listed in `skippedMethods` bypasses authentication, either
by exact full method name or by a service prefix ending in `/`.

Handlers can then read the caller identity:

```go
subject, ok := auth.SubjectFromContext(ctx)
```

## Wiring the client

```go
conn, err := grpc.NewClient(addr,
    grpc.WithTransportCredentials(credentials),
    grpc.WithUnaryInterceptor(auth.NewUnaryClientInterceptor(tm, "scheduler", time.Minute)),
    grpc.WithStreamInterceptor(auth.NewStreamClientInterceptor(tm, "scheduler", time.Minute)),
)
```

A token can also be attached manually with
`auth.AppendTokenToOutgoingContext(ctx, token)`.

## Behavior

| Case | Result |
| --- | --- |
| No `authorization` metadata | `codes.Unauthenticated` |
| Non-bearer scheme / empty token | `codes.Unauthenticated` |
| Expired or not-yet-valid token | `codes.Unauthenticated` |
| Signature signed with another key | `codes.Unauthenticated` |
| Algorithm other than HS256 (including `none`) | `codes.Unauthenticated` |
| Issuer (or configured audience) mismatch | `codes.Unauthenticated` |
| Token without `exp` | `codes.Unauthenticated` |
| Valid token | Request reaches the handler |

Token TTL is controlled by the caller of `IssueToken`; a non-positive TTL falls
back to `auth.DefaultTokenTTL` (10 minutes). Prefer short-lived tokens and
re-issue them per call or per connection.

## Tests

Tests start a real gRPC server on `bufconn` with the interceptors installed and
use the standard gRPC health service as the target:

```bash
go test ./pkg/rpc/auth/...
```

Covered cases: missing token, expired token, wrong signing key, wrong issuer,
`none` algorithm, malformed header, valid token (unary and stream), health
service exemption, and claim propagation to handlers.
