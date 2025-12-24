# Dragonfly Injector Deployment Guide

This guide provides step-by-step instructions for deploying the Dragonfly Injector webhook to your Kubernetes cluster.

## Prerequisites

- Kubernetes cluster v1.11.3 or higher
- `kubectl` configured to access your cluster
- [cert-manager](https://cert-manager.io/) installed in your cluster
- Cluster admin permissions

## Quick Start

### 1. Install cert-manager (if not already installed)

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

Wait for cert-manager to be ready:

```bash
kubectl wait --for=condition=Available --timeout=300s deployment/cert-manager -n cert-manager
kubectl wait --for=condition=Available --timeout=300s deployment/cert-manager-webhook -n cert-manager
```

### 2. Deploy the Injector

Using kustomize (recommended):

```bash
kubectl apply -k deploy/injector/
```

Or apply manifests individually:

```bash
kubectl apply -f deploy/injector/namespace.yaml
kubectl apply -f deploy/injector/rbac.yaml
kubectl apply -f deploy/injector/configmap.yaml
kubectl apply -f deploy/injector/certificate.yaml
kubectl apply -f deploy/injector/service.yaml
kubectl apply -f deploy/injector/deployment.yaml
kubectl apply -f deploy/injector/webhook.yaml
```

### 3. Verify Installation

Check that all resources are created:

```bash
kubectl get all -n dragonfly-system -l component=injector
```

Check webhook configuration:

```bash
kubectl get mutatingwebhookconfiguration dragonfly-injector-webhook
```

Check certificate status:

```bash
kubectl get certificate -n dragonfly-system dragonfly-injector-serving-cert
```

Wait for the injector pods to be ready:

```bash
kubectl wait --for=condition=Ready --timeout=300s pod -n dragonfly-system -l component=injector
```

## Configuration

### Customize Injector Configuration

Edit the ConfigMap to customize injector behavior:

```bash
kubectl edit configmap dragonfly-injector-config -n dragonfly-system
```

Available configuration options:

```yaml
enable: true                                    # Enable/disable injection
proxy_port: 4001                               # Dfdaemon proxy port
cli_tools_image: "dragonflyoss/toolkits:latest"  # CLI tools image
cli_tools_dir_path: "/dragonfly-tools"         # Tools directory
```

Changes are automatically reloaded within 15 seconds.

### Customize Image Version

Edit the kustomization.yaml to use a specific version:

```yaml
images:
  - name: dragonflyoss/injector
    newTag: v2.3.5  # Specify your version
```

Then apply:

```bash
kubectl apply -k deploy/injector/
```

### Adjust Resource Limits

Edit the deployment to adjust resource requests/limits:

```bash
kubectl edit deployment dragonfly-injector -n dragonfly-system
```

## Usage

### Enable Injection for a Namespace

Add the label to enable injection for all pods in a namespace:

```bash
kubectl label namespace my-namespace dragonfly.io/inject=true
```

### Enable Injection for a Specific Pod

Add the annotation to a pod:

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

### Verify Injection

Create a test pod and verify injection:

```bash
# Create test namespace with injection enabled
kubectl create namespace test-injection
kubectl label namespace test-injection dragonfly.io/inject=true

# Create test pod
kubectl run test-pod --image=nginx -n test-injection

# Check if injection occurred
kubectl get pod test-pod -n test-injection -o yaml | grep -A 10 "dragonfly"
```

You should see:
- Environment variables (NODE_NAME, DRAGONFLY_PROXY_PORT, DRAGONFLY_INJECT_PROXY, DRAGONFLY_TOOLS_PATH)
- Init container (d7y-cli-tools)
- Volume mounts (dfdaemon-unix-sock, d7y-cli-tools-volume)

## Troubleshooting

### Check Injector Logs

```bash
kubectl logs -n dragonfly-system -l component=injector --tail=100 -f
```

### Check Webhook Status

```bash
kubectl describe mutatingwebhookconfiguration dragonfly-injector-webhook
```

### Check Certificate

```bash
kubectl get certificate -n dragonfly-system dragonfly-injector-serving-cert
kubectl describe certificate -n dragonfly-system dragonfly-injector-serving-cert
```

### Common Issues

#### 1. Pods Not Being Injected

**Symptoms**: Pods created but no dragonfly resources injected

**Solutions**:
- Verify namespace has the label: `kubectl get namespace <namespace> --show-labels`
- Check webhook configuration: `kubectl get mutatingwebhookconfiguration`
- Check injector logs for errors
- Verify injector pods are running: `kubectl get pods -n dragonfly-system`

#### 2. Certificate Issues

**Symptoms**: Webhook fails with TLS/certificate errors

**Solutions**:
- Verify cert-manager is running: `kubectl get pods -n cert-manager`
- Check certificate status: `kubectl get certificate -n dragonfly-system`
- Check certificate secret exists: `kubectl get secret dragonfly-injector-webhook-cert -n dragonfly-system`
- Recreate certificate: `kubectl delete certificate dragonfly-injector-serving-cert -n dragonfly-system`

#### 3. Webhook Timeout

**Symptoms**: Pod creation times out or fails

**Solutions**:
- Check injector pod health: `kubectl get pods -n dragonfly-system`
- Increase webhook timeout in webhook.yaml
- Check network policies aren't blocking webhook traffic
- Verify service endpoints: `kubectl get endpoints -n dragonfly-system`

#### 4. Leader Election Issues

**Symptoms**: Multiple injectors trying to be leader

**Solutions**:
- Check RBAC permissions for leases
- Check injector logs for leader election messages
- Verify only one leader is active

## Monitoring

### Metrics

The injector exposes Prometheus metrics on port 8443:

```bash
kubectl port-forward -n dragonfly-system svc/dragonfly-injector-webhook-service 8443:8443
```

Then access metrics at `https://localhost:8443/metrics`

### Health Checks

Health endpoint:
```bash
kubectl port-forward -n dragonfly-system deployment/dragonfly-injector 8081:8081
curl http://localhost:8081/healthz
```

Readiness endpoint:
```bash
curl http://localhost:8081/readyz
```

## Upgrading

### Upgrade Injector Version

1. Update the image tag in kustomization.yaml or deployment
2. Apply the changes:
   ```bash
   kubectl apply -k deploy/injector/
   ```
3. Wait for rollout to complete:
   ```bash
   kubectl rollout status deployment/dragonfly-injector -n dragonfly-system
   ```

### Rollback

If issues occur, rollback to previous version:

```bash
kubectl rollout undo deployment/dragonfly-injector -n dragonfly-system
```

## Uninstallation

Remove all injector resources:

```bash
kubectl delete -k deploy/injector/
```

Or remove individually:

```bash
kubectl delete mutatingwebhookconfiguration dragonfly-injector-webhook
kubectl delete namespace dragonfly-system
```

## Security Considerations

- The injector runs as non-root user (UID 65532)
- Read-only root filesystem
- HTTP/2 disabled by default
- TLS required for all webhook communication
- RBAC permissions limited to reading namespaces and pods
- Secrets managed by cert-manager

## Advanced Configuration

### Custom TLS Certificates

To use custom certificates instead of cert-manager:

1. Create a secret with your certificates:
   ```bash
   kubectl create secret tls dragonfly-injector-webhook-cert \
     --cert=path/to/tls.crt \
     --key=path/to/tls.key \
     -n dragonfly-system
   ```

2. Remove the certificate.yaml from deployment
3. Update the MutatingWebhookConfiguration with your CA bundle

### High Availability

For production deployments:

1. Increase replica count in deployment.yaml
2. Enable leader election (already enabled by default)
3. Configure pod disruption budget:
   ```yaml
   apiVersion: policy/v1
   kind: PodDisruptionBudget
   metadata:
     name: dragonfly-injector-pdb
     namespace: dragonfly-system
   spec:
     minAvailable: 1
     selector:
       matchLabels:
         app: dragonfly
         component: injector
   ```

### Network Policies

Restrict injector network access:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: dragonfly-injector
  namespace: dragonfly-system
spec:
  podSelector:
    matchLabels:
      component: injector
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector: {}
      ports:
        - protocol: TCP
          port: 9443
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
      ports:
        - protocol: TCP
          port: 443
```

## Support

For issues and questions:
- GitHub Issues: https://github.com/dragonflyoss/dragonfly/issues
- Slack: #dragonfly on CNCF Slack
- Documentation: https://d7y.io/
