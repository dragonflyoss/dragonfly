# Injector main-repository migration

The Pod injector was migrated from `dragonflyoss/dragonfly-injector` at commit
`6b54b9bd` into the Dragonfly main module. The production code is in
`internal/injector`, the Pod admission webhook is in
`internal/injector/webhook/v1`, and the executable entry point is
`cmd/injector/main.go`.

This increment preserves the baseline injection behavior: binary and dfdaemon
Unix-socket injection, namespace-label or Pod-annotation opt-in, file-backed
configuration, and periodic configuration reload. Its unit tests use fake
Kubernetes clients and do not require envtest binaries.

## Build

Build the host binary with:

```console
make build-injector
```

The binary is written to `bin/<goos>_<goarch>/injector`. The container targets
follow the existing Dragonfly image conventions and honor `D7Y_REGISTRY` and
`D7Y_VERSION`:

```console
make docker-build-injector
make docker-push-injector
```

The image is defined in `build/images/injector/Dockerfile` and runs
`/opt/dragonfly/bin/injector`. No image is published by the source migration
itself.

## Deploy

The focused Kustomize manifests are in `config/injector`. They create the
`dragonfly-system` namespace, service account and RBAC, configuration,
cert-manager issuer and certificate, Deployment, webhook Service, and
MutatingWebhookConfiguration. Cert-manager must already be installed to issue
the serving certificate and inject its CA bundle. Build or otherwise provide
`dragonflyoss/injector:latest`, or override the image in a local Kustomize
overlay, then run:

```console
kubectl apply -k config/injector
```

Enable injection for a namespace with the `dragonfly.io/inject=true` label, or
for an individual Pod with the `dragonfly.io/inject=true` annotation. A Pod
annotation takes precedence over the namespace label. Set
`dragonfly.io/skip-unix-sock-inject=true` to skip only Unix-socket injection,
and `dragonfly.io/init-container-image` to override the binaries init-container
image for a Pod.

## Configuration

The Deployment mounts `dragonfly-injector-config` at `/etc/dragonfly`; the
injector reads `/etc/dragonfly/injector.yaml` and checks it for changes every
30 seconds. The `initContainerImage` fields are `registry`, `repository`, `tag`,
`digest`, `pullPolicy`, and `pullSecrets`. Invalid or missing files fall back to
the compiled defaults.

The command also exposes the baseline controller-runtime options, including
`--webhook-cert-path`, `--webhook-cert-name`, `--webhook-cert-key`,
`--metrics-bind-address`, `--metrics-secure`, metrics certificate flags,
`--health-probe-bind-address`, `--leader-elect`, and `--enable-http2`. Metrics
are disabled in the supplied Deployment. Leader election requires additional
Lease RBAC and an appropriate multi-replica deployment configuration.

## Validation boundary

Focused formatting, unit-test, and build checks belong to this migration. The
manifests are deployment assets, not evidence of a live-cluster test. Helm
chart integration and corrections remain in their respective chart changes;
this document does not claim chart work, image publication, envtest, a live
Kubernetes deployment, or full end-to-end validation.
