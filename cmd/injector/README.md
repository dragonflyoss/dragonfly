# Dragonfly Injector

A Kubernetes Mutating Admission Webhook that automatically injects Dragonfly P2P capabilities into pods.

## Overview

The Dragonfly Injector simplifies Kubernetes pod configuration by automating the injection of Dragonfly's P2P proxy settings, dfdaemon socket mounts, and CLI tools through annotation-based policies.

## Features

### 1. Annotation/Label-Based Injection

The webhook supports injecting P2P configurations based on annotations at the pod level or labels at the namespace level.

**Priority**: Pod-level annotations > Namespace-level labels

#### Namespace-Level Injection

Enable injection for all pods in a namespace:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    dragonfly.io/inject: "true"
  name: my-namespace
```

#### Pod-Level Injection

Enable injection for a specific pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
  annotations:
    dragonfly.io/inject: "true"
spec:
  containers:
    - name: app
      image: my-app:latest
```

### 2. P2P Proxy Environment Variables

Automatically injects environment variables to route application traffic through the Dragonfly P2P proxy:

- `NODE_NAME`: Node where pod is scheduled (via Downward API)
- `DRAGONFLY_PROXY_PORT`: Dfdaemon proxy port (default: 4001)
- `DRAGONFLY_INJECT_PROXY`: Constructed proxy URL (`http://$(NODE_NAME):$(DRAGONFLY_PROXY_PORT)`)

### 3. Dfdaemon Socket Volume Mount

Mounts the dfdaemon Unix socket to enable CLI tools (dfget, etc.) to communicate with the dfdaemon daemon:

- **Volume**: hostPath volume for `/var/run/dragonfly/dfdaemon.sock`
- **Mount**: Mounted in all containers at the same path

### 4. CLI Tools Injection

Injects Dragonfly CLI tools (dfget, etc.) via init container for base images that don't include them:

- **Init Container**: Copies tools from multi-arch toolkit image
- **Shared Volume**: emptyDir volume for tool binaries
- **Environment Variable**: `DRAGONFLY_TOOLS_PATH` points to tools directory

#### Custom Tools Image

Specify a custom CLI tools image via annotation:

```yaml
metadata:
  annotations:
    dragonfly.io/inject: "true"
    dragonfly.io/cli-tools-image: "my-registry/dragonfly-tools:v1.0.0"
```

## Configuration

The injector is configured via a ConfigMap mounted at `/etc/dragonfly-injector/config.yaml`:

```yaml
enable: true                                    # Enable/disable injection
proxy_port: 4001                               # Dfdaemon proxy port
cli_tools_image: "dragonflyoss/toolkits:latest"  # CLI tools image
cli_tools_dir_path: "/dragonfly-tools"         # Tools directory path
```

The configuration is hot-reloaded every 15 seconds.

## Deployment

### Prerequisites

- Kubernetes v1.11.3+
- cert-manager (for automatic TLS certificate management)

### Build

```bash
# Build binary
make build-injector

# Build Docker image
make docker-build-injector

# Push Docker image
D7Y_REGISTRY=myregistry D7Y_VERSION=v1.0.0 make docker-push-injector
```

### Install

The injector requires:

1. **TLS Certificates**: Managed by cert-manager or provided manually
2. **MutatingWebhookConfiguration**: Registers the webhook with Kubernetes
3. **RBAC**: ServiceAccount, Role, and RoleBinding for namespace/pod access
4. **Deployment**: Runs the injector webhook server

## Command-Line Flags

```bash
dragonfly-injector [flags]
```

### Flags

- `--metrics-bind-address`: Metrics endpoint address (default: "0")
- `--health-probe-bind-address`: Health probe address (default: ":8081")
- `--leader-elect`: Enable leader election (default: false)
- `--metrics-secure`: Serve metrics securely via HTTPS (default: true)
- `--webhook-cert-path`: Directory containing webhook TLS certificate
- `--webhook-cert-name`: Webhook certificate filename (default: "tls.crt")
- `--webhook-cert-key`: Webhook key filename (default: "tls.key")
- `--metrics-cert-path`: Directory containing metrics server certificate
- `--metrics-cert-name`: Metrics certificate filename (default: "tls.crt")
- `--metrics-cert-key`: Metrics key filename (default: "tls.key")
- `--enable-http2`: Enable HTTP/2 (default: false, disabled for security)

## Example

Complete example of an injected pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
  annotations:
    dragonfly.io/inject: "true"
spec:
  initContainers:
    - name: d7y-cli-tools
      image: dragonflyoss/toolkits:latest
      command: ["cp", "-rf", "/dragonfly-tools/.", "/dragonfly-tools-mount/"]
      volumeMounts:
        - name: d7y-cli-tools-volume
          mountPath: /dragonfly-tools-mount
  containers:
    - name: app
      image: my-app:latest
      env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: DRAGONFLY_PROXY_PORT
          value: "4001"
        - name: DRAGONFLY_INJECT_PROXY
          value: "http://$(NODE_NAME):$(DRAGONFLY_PROXY_PORT)"
        - name: DRAGONFLY_TOOLS_PATH
          value: "/dragonfly-tools-mount"
      volumeMounts:
        - name: dfdaemon-unix-sock
          mountPath: /var/run/dragonfly/dfdaemon.sock
        - name: d7y-cli-tools-volume
          mountPath: /dragonfly-tools-mount
  volumes:
    - name: dfdaemon-unix-sock
      hostPath:
        path: /var/run/dragonfly/dfdaemon.sock
        type: Socket
    - name: d7y-cli-tools-volume
      emptyDir: {}
```

## Security

- Runs as non-root user (UID 65532)
- HTTP/2 disabled by default to prevent CVE vulnerabilities
- TLS required for webhook communication
- RBAC-based access control

## Troubleshooting

### Check Webhook Status

```bash
kubectl get mutatingwebhookconfigurations
kubectl describe mutatingwebhookconfiguration dragonfly-injector
```

### View Injector Logs

```bash
kubectl logs -n dragonfly-system deployment/dragonfly-injector
```

### Verify Injection

Check if a pod has been injected:

```bash
kubectl get pod <pod-name> -o yaml | grep -A 5 "dragonfly"
```

### Common Issues

1. **Injection not working**: Verify annotation/label is set correctly
2. **Certificate errors**: Check cert-manager is running and certificates are valid
3. **Permission errors**: Verify RBAC permissions for namespace/pod access

## License

Copyright 2025 The Dragonfly Authors.

Licensed under the Apache License, Version 2.0.
