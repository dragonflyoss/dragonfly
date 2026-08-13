# JWT Authentication for Inter-Component gRPC Requests

Status: Proposed

Tracking issue: <https://github.com/dragonflyoss/dragonfly/issues/4417>

## 1. Summary

Dragonfly components currently trust any caller that can reach their gRPC
business endpoints unless mTLS is enabled. This proposal adds token-based
authentication to remote inter-component gRPC requests by using short-lived
JWTs in gRPC metadata.

The first version deliberately uses a cluster-wide shared HMAC key. Every
component that initiates a remote gRPC request signs a JWT locally, and every
component that accepts a remote business request verifies the JWT locally.
Manager is not a token issuer and no token distribution or refresh RPC is
introduced.

The authentication layer is independent of the existing transport security
implementation. This proposal does not change the current TLS or mTLS code.

## 2. Goals

- Authenticate remote inter-component gRPC business requests.
- Use the same protocol and validation rules in Go and Rust.
- Support rolling upgrades from unauthenticated deployments.
- Support key rotation without a cluster-wide authentication outage.
- Keep the existing protobuf APIs unchanged.
- Keep Manager REST user authentication separate from internal gRPC
  authentication.
- Fail closed when authentication is enabled but its configuration is invalid.

## 3. Non-goals

- Identifying or authorizing individual components.
- Restricting particular components to particular RPC methods.
- Adding a Manager token issuance, refresh, or revocation service.
- Integrating an external OAuth, OIDC, SPIFFE, or identity provider.
- Replacing or modifying TLS or mTLS.
- Authenticating local dfdaemon Download RPCs over Unix domain sockets.
- Detecting replay of a valid JWT within its lifetime.
- Reloading keys without restarting a component in the first version.
- Re-authenticating individual messages within an existing gRPC stream.

## 4. Security model

Possession of any active cluster key proves only that the caller is a member of
the Dragonfly deployment. All holders of the shared key can mint a valid token
for any Dragonfly audience. Claims such as a component name or role therefore
must not be used for authorization in the first version.

Compromise of one dfdaemon or Seed Peer that can read the shared key compromises
the JWT authentication boundary for the whole cluster. This blast radius is an
accepted first-version tradeoff. Per-component asymmetric credentials can be
introduced later without changing the gRPC metadata transport.

JWTs are bearer credentials and do not encrypt traffic. When the underlying
channel is plaintext, this feature protects against callers that can reach a
gRPC port but cannot read traffic in transit. It does not protect against an
on-path attacker. Transport security remains a separate deployment decision.

## 5. Protocol contract

### 5.1 Metadata

The client sends exactly one lowercase gRPC metadata entry:

```text
authorization: Bearer <compact-jwt>
```

The server rejects multiple authorization values, an empty value, an unknown
scheme, a malformed JWT, and credentials larger than 4 KiB.

### 5.2 JOSE header

```json
{
  "alg": "HS256",
  "typ": "dragonfly-grpc+jwt",
  "kid": "key-2026-08"
}
```

- `alg` is fixed to `HS256`. The validator must not select an algorithm from
  untrusted token input.
- `typ` is fixed to `dragonfly-grpc+jwt` and separates this token from Manager
  REST user JWTs and other JWT types.
- `kid` selects a key from the local trusted keyring.

### 5.3 Claims

```json
{
  "iss": "dragonfly-internal",
  "aud": "urn:dragonfly:grpc:scheduler",
  "iat": 1786435200,
  "exp": 1786435800
}
```

Required claims:

- `iss`: exact configured issuer.
- `aud`: exact audience expected by the target server.
- `iat`: token creation time.
- `exp`: token expiration time.

The first version does not emit or trust `sub`, role, scope, component identity,
or permission claims. `jti` is omitted because no replay cache is maintained.

Audience values are constants in both implementations:

| Target service | Audience |
| --- | --- |
| Manager | `urn:dragonfly:grpc:manager` |
| Scheduler | `urn:dragonfly:grpc:scheduler` |
| dfdaemon | `urn:dragonfly:grpc:dfdaemon` |

Audience selection is not configurable at individual call sites.

### 5.4 Validation

The server requires all of the following:

1. Exactly one Bearer credential is present for protected methods.
2. `alg`, `typ`, `kid`, signature, issuer, and audience are valid.
3. `iat` is not later than the server time plus the configured clock skew.
4. `exp` is later than the server time minus the configured clock skew.
5. `exp - iat` is positive and no greater than `maxTokenTTL`.

Authentication failures return `codes.Unauthenticated` with a generic message.
Detailed failure reasons are recorded only in bounded metrics. Raw tokens,
authorization metadata, secrets, and full decoded claims must never be logged.

## 6. Configuration contract

Manager, Scheduler, and dfdaemon use the same top-level configuration shape:

```yaml
grpcAuth:
  # disabled, permissive, or required
  mode: disabled

  # If true, Go PerRPCCredentials and Rust client construction reject a
  # plaintext transport. It is false by default for compatibility with
  # existing Dragonfly deployments.
  requireTransportSecurity: false

  jwt:
    issuer: dragonfly-internal
    tokenTTL: 10m
    maxTokenTTL: 15m
    clockSkew: 30s
    refreshBefore: 1m

    activeKeyID: key-2026-08
    keys:
      - id: key-2026-08
        secretFile: /etc/dragonfly/secrets/grpc-jwt/key-2026-08
      - id: key-2026-07
        secretFile: /etc/dragonfly/secrets/grpc-jwt/key-2026-07
```

Modes have symmetric client and server behavior:

| Mode | Client behavior | Server behavior |
| --- | --- | --- |
| `disabled` | Does not add a token | Does not authenticate |
| `permissive` | Adds a token | Accepts a missing token, rejects an invalid supplied token |
| `required` | Adds a token | Rejects a missing or invalid token |

`permissive` is a migration state, not a secure steady state. Because a missing
token is accepted, it is vulnerable to metadata stripping.

### 6.1 Key file format

- `secretFile` contains one RFC 4648 standard Base64 string.
- After trimming leading and trailing whitespace, the content is decoded
  as Base64.
- The decoded key must contain at least 32 bytes. Production keys must be
  generated by a cryptographically secure random generator.
- The path must refer to a regular readable file.
- Duplicate key IDs and a missing `activeKeyID` are configuration errors.
- When mode is `permissive` or `required`, missing or invalid keys are fatal
  startup errors.
- Keys are read once during startup and retained only in memory.
- Production deployments should mount files with mode `0400` or `0440`.

A suitable key can be created with:

```bash
umask 077
openssl rand -base64 32 > key-2026-08
```

All components in one cluster receive the same keyring. The configuration
contains file paths, not inline secret values. Manager does not distribute
these files at runtime; Docker Compose, Kubernetes, or the host provisioning
system does so before process startup.

Manager's existing `auth.jwt` configuration continues to authenticate REST
users. It must use a different secret and a separate validator. A Manager REST
JWT must never authenticate an internal gRPC request.

### 6.2 Configuration validation

Configuration loading rejects:

- An unknown mode.
- A missing issuer in an enabled mode.
- `tokenTTL < 1s`, `maxTokenTTL < 1s`, or `tokenTTL > maxTokenTTL`.
- A JWT duration that is not an exact number of seconds.
- `clockSkew < 0`.
- `refreshBefore <= 0` or `refreshBefore >= tokenTTL`.
- A missing, duplicate, unreadable, invalid, or short key.
- An `activeKeyID` that is absent from the keyring.

Partially configured authentication must never silently fall back to disabled.

## 7. Request lifecycle

Each process constructs one concurrency-safe token provider during startup. The
provider caches one token per `(audience, activeKeyID)` pair.

- A cached token is reused while more than `refreshBefore` remains.
- A new token is signed locally when the cached token approaches expiration.
- Signing or metadata construction failure fails the RPC locally. The client
  must not retry without authentication.
- Manager is not contacted during token generation or refresh.

For streaming RPCs, authentication happens when the stream is established. A
stream that was authenticated successfully continues after the token expires.
A reconnect obtains a new token. Mode or key changes in the first version
require a process restart, which also terminates old streams.

## 8. Protected call matrix

| Caller | Target | Repository | Required work |
| --- | --- | --- | --- |
| Scheduler | Manager | `dragonfly` | Add Manager audience credentials |
| dfdaemon / Seed Peer | Manager | `client` | Add Manager audience interceptor |
| dfdaemon / Seed Peer | Scheduler | `client` | Add Scheduler audience interceptor |
| Scheduler | dfdaemon / Seed Peer | `dragonfly` | Centralize all peer dial options |
| dfdaemon | dfdaemon gRPC services | `client` | Add dfdaemon audience interceptor |
| `dfctl task preheat` | Scheduler | `client` | Load gRPC auth config and add a token |

Manager v1/v2, Scheduler v1/v2, and dfdaemon business methods are
protected. The gRPC health service is always public so that liveness and
readiness probes do not need the cluster secret. Reflection is protected in
`required` mode.

Local dfdaemon Download RPCs over Unix domain sockets are outside this proposal.

Old clients and legacy services that cannot attach or validate JWTs work only
while the relevant server is `disabled` or `permissive`.

## 9. Go implementation

Add a transport-independent package under `pkg/rpc/auth/jwt` with:

- Configuration and keyring loading.
- Claims and protocol constants.
- A concurrency-safe token provider implementing
  `credentials.PerRPCCredentials`.
- A verifier.
- Unary and stream server interceptors.
- Authentication metrics and bounded failure reasons.

`RequireTransportSecurity()` returns the configured
`requireTransportSecurity` value. No TLS or mTLS implementation is changed.

Manager and Scheduler add the authentication interceptor to both existing
unary and stream interceptor chains. Authentication runs before request
validation and business handlers. The health service is explicitly bypassed;
reflection and business methods are not.

Inject a single authentication dial option into each existing dial-options
assembly path. The authentication option owns:

- JWT Per-RPC credentials and target audience.

Existing transport credentials and OpenTelemetry setup remain unchanged and
are composed with the authentication option at process assembly time. All
standard, Seed Peer, job, persistent-task, and persistent-cache-task paths must
receive the authentication option. In particular, the persistent-task paths in
`scheduler/service/service_v2.go` currently construct clients locally, so the
auth option must be injected into the V2 service and appended there instead of
being rediscovered from global configuration. This proposal does not change
their TLS or mTLS behavior.

The JWT library is declared as a direct Go module dependency. Manager's REST
authentication package is not reused.

## 10. Rust implementation

Add matching configuration types to `dragonfly-client-config` and a dedicated
`dragonfly-client-auth` workspace crate with:

- Keyring loading and strict validation.
- Matching claims, JOSE header, audiences, and clock rules.
- A concurrency-safe token provider.
- Client and server Tonic interceptors.
- Matching bounded metrics and error behavior.

JWT injection is composed with the existing tracing interceptor so both tracing
and authorization metadata are retained.

Update:

- Manager client.
- Scheduler client.
- dfdaemon DfdaemonUpload gRPC clients. The maintained Rust client no longer
  contains the legacy Seeder service.
- dfdaemon business gRPC servers.
- Reflection service authentication.
- `dfctl task preheat` direct Scheduler calls.

The health service remains outside the auth interceptor. Local dfdaemon Download
clients and servers over Unix domain sockets are unchanged.

The JWT crate is declared as a direct dependency at a version compatible with
the repository's Rust toolchain. The implementation must not rely on a
transitive `jsonwebtoken` dependency from `Cargo.lock`.

## 11. Observability

Expose a bounded server counter equivalent to:

```text
dragonfly_grpc_auth_requests_total{audience,mode,result,reason}
```

Allowed `reason` values are:

- `none`
- `missing`
- `malformed`
- `unsupported_alg`
- `invalid_type`
- `unknown_kid`
- `invalid_signature`
- `invalid_issuer`
- `invalid_audience`
- `expired`
- `invalid_iat`
- `ttl_exceeded`

Also expose token generation, generation failure, cache hit, and cache miss
counters. RPC methods and remote addresses are not metric labels because both
can be attacker-controlled. Remote addresses may be included only in sampled
logs.

Startup logs may report the mode, active key ID, and trusted key count. They
must not report keys or tokens.

## 12. Rolling upgrade

When mode is omitted or `disabled`, no keyring is required, clients do not add
credentials, and servers do not authenticate requests. A normal application or
chart rolling upgrade therefore preserves the legacy behavior.

The first authenticated rollout uses one global mode transition:

1. Create and distribute the shared keyring, deploy JWT-capable binaries, and
   set the global mode to `permissive`. These changes may be combined in one
   rolling update: upgraded clients send tokens, upgraded servers accept both
   authenticated requests and requests from old callers, and old servers
   ignore the additional authorization metadata.
2. Wait for every component and maintained external caller, including
   `dfctl task preheat`, to be upgraded. Monitor `missing` authentication
   metrics until every expected business call path sends tokens.
3. Change the global mode to `required`. This transition is safe as a rolling
   update because clients in both `permissive` and `required` modes send tokens,
   while servers in both modes accept valid tokens.

Once a deployment is in `required` mode, later application upgrades keep that
mode and use the normal one-step rolling update as long as the JWT protocol and
trusted keyring remain compatible.

Each mode change restarts the process in the first version, ensuring that
streams established before `required` mode do not remain unauthenticated.

Rollback changes the global mode from `required` to `permissive` before
disabling authentication. Older servers ignore the extra authorization
metadata sent by upgraded clients.

## 13. Key rotation

1. Distribute a new key to every component as a trusted key while the old key
   remains active.
2. Restart all components so every verifier accepts both keys.
3. Change `activeKeyID` to the new key and restart callers incrementally.
4. Wait at least `maxTokenTTL + clockSkew` after the last caller switches.
5. Remove the old key and restart servers. The restart also terminates streams
   authenticated with the old key.

## 14. Testing

### 14.1 Unit tests in both languages

- Valid token signing and validation.
- Invalid algorithm, type, key ID, signature, issuer, and audience.
- Missing and duplicate authorization metadata.
- Expired, future-issued, and excessive-lifetime tokens.
- All three modes.
- Token cache refresh behavior under concurrency.
- Unary and stream interceptor behavior.
- Health bypass and reflection protection.
- Strict configuration and key-file validation.

### 14.2 Cross-language contract tests

Commit the same deterministic test vector to both repositories with a fixed
test key, header, claims, and clock:

- Go signs and Rust validates.
- Rust signs and Go validates.
- Both produce and accept the same audience and time semantics.

The fixture keeps one deterministic token from each implementation. Their
compact strings may differ because JSON member order in the JOSE header is not
semantically significant; cross-language signature validation is the contract.

The fixture secret is test-only and clearly marked as unusable in production.

### 14.3 Integration tests

- Scheduler to Manager.
- Rust dfdaemon to Manager and Scheduler.
- Scheduler standard, persistent, and persistent-cache callbacks to dfdaemon.
- Rust peer-to-peer Upload calls.
- `dfctl task preheat` to Scheduler.
- A stream starts with a valid token, survives token expiration, and reconnects
  with a refreshed token.
- Mixed old/new binaries in `permissive` mode.
- Missing and invalid credentials in `required` mode.
- Overlapping old/new keys during rotation.
- Existing Manager REST authentication remains unchanged.
- Local Unix domain socket calls and health probes remain unchanged.

## 15. Repository and merge order

Three repositories require changes:

1. `dragonflyoss/dragonfly`
2. `dragonflyoss/client`
3. `dragonflyoss/helm-charts`

`dragonflyoss/api` is not changed because authorization is transported in gRPC
metadata and no protobuf message or token RPC is introduced.

Recommended merge sequence:

1. **dragonfly PR 1: protocol and Go implementation**
   - Land this design and protocol constants.
   - Add Go keyring, provider, verifier, interceptors, config, unit tests, and
     authentication metrics.
   - Protect Manager and Scheduler business servers.
   - Add credentials to every Scheduler outbound client path without changing
     the existing TLS or mTLS behavior.
   - Keep the default mode `disabled`.
2. **client PR: Rust implementation**
   - Implement the frozen protocol and matching test vector.
   - Update Manager, Scheduler, Upload, and `dfctl` call paths.
   - Protect the Upload business server.
   - Keep the default mode `disabled`.
3. **dragonfly PR 2: integration**
   - Advance the `client` submodule to the merged client commit.
   - Add cross-component compatibility and end-to-end tests.
   - Add Docker Compose examples and migration documentation.
4. **helm-charts PR: deployment support**
   - Add one `existingSecret` reference for the cluster JWT keyring.
   - Mount the same keyring into Manager, Scheduler, dfdaemon, and Seed Peer.
   - Render matching `grpcAuth` configuration for every workload.
   - Do not put secret values directly into ordinary chart values.
5. **dragonfly PR 3: chart pointer**
   - Advance the `deploy/helm-charts` submodule after the chart PR merges.
   - Run the final mixed-version and required-mode E2E suites.

The Go and Rust implementation PRs may be developed concurrently after this
contract is frozen. The first authenticated rollout uses `permissive`, and a
deployment must not enable `required` until every component and maintained
external caller runs a compatible version and sends tokens.

## 16. Acceptance criteria

- In `required` mode, every remote business RPC without a valid JWT returns
  `Unauthenticated` before entering its business handler.
- Every maintained Go and Rust remote client path sends the correct audience.
- Manager REST JWTs are rejected by internal gRPC validators.
- Health checks and local Unix domain socket RPCs continue to work without JWT.
- An omitted or `disabled` configuration preserves legacy behavior and requires
  no key files.
- Mixed-version rolling upgrades work in `permissive` mode.
- The global `permissive` to `required` transition works as a rolling upgrade.
- Key rotation works with an overlap window and without authentication outage.
- Tokens and secrets never appear in logs, errors, metrics labels, or traces.
- No protobuf API change or Manager token service is introduced.
