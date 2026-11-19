# Dragonfly Data Controller

The data controller is a Kubernetes controller that manages data in the Dragonfly P2P network
declaratively. Instead of calling the manager's preheat, get task and delete task jobs by hand,
you declare a `Dataset` and the controller keeps the P2P network in the declared state for the
whole life of the data: it distributes the data, keeps it warm, verifies that peers still hold it,
and removes it when it expires or when the `Dataset` is deleted.

## Concepts

### Dataset

A `Dataset` is a namespaced resource that declares a set of data sources and how they are managed.

- **Sources** are the data that make up the dataset. A source is either a `file` source (one or more
  URLs served over HTTP(S), object storage or HDFS) or an `image` source (an OCI manifest URL, whose
  layers are all managed). Every source has a name that the status reports on.
- **Distribution** controls how the data is pushed into the P2P network: the scope
  (`single_seed_peer`, `all_seed_peers` or `all_peers`), the peers to target (by IP, count or
  percentage), the concurrency, the scheduler clusters, the job timeout and the retry limit.
- **Lifecycle** controls what happens over time: a TTL after which the data expires, an expire action
  (`Purge` or `Retain`), a refresh interval or cron schedule that re-distributes the data, a
  verification interval that checks which peers hold the data, and a deletion policy (`Purge` or
  `Retain`) applied when the `Dataset` is deleted.
- **ManagerRef** optionally points the dataset at a manager other than the one configured on the
  controller, with the personal access token read from a `Secret`.

Credentials never appear in the spec. Basic authentication, extra HTTP headers, object storage
credentials, HDFS delegation tokens and the manager token are all read from `Secret` references in
the same namespace.

```yaml
apiVersion: data.d7y.io/v1alpha1
kind: Dataset
metadata:
  name: llm-weights
  labels:
    tier: hot
spec:
  sources:
    - name: weights
      type: file
      urls:
        - https://models.example.com/llama/model-00001-of-00002.safetensors
        - https://models.example.com/llama/model-00002-of-00002.safetensors
      tag: v3
  distribution:
    scope: all_seed_peers
  lifecycle:
    refreshInterval: 6h
    verifyInterval: 1h
    ttl: 48h
    deletionPolicy: Purge
```

### DataLifecyclePolicy

A `DataLifecyclePolicy` applies a lifecycle to every `Dataset` in its namespace that matches its label
selector, so that lifecycle rules are managed centrally instead of being repeated on every dataset.
When several policies match a dataset, the one with the highest `priority` wins, ties are broken by
name. Fields set in the dataset's own `spec.lifecycle` override the policy. Refresh settings are taken
as a pair: a dataset that sets `refreshInterval` or `refreshSchedule` ignores the policy's refresh
settings entirely.

```yaml
apiVersion: data.d7y.io/v1alpha1
kind: DataLifecyclePolicy
metadata:
  name: hot-data
spec:
  datasetSelector:
    matchLabels:
      tier: hot
  priority: 10
  lifecycle:
    refreshSchedule: "0 2 * * *"
    verifyInterval: 2h
    deletionPolicy: Purge
```

The policy status reports how many datasets it governs and whether it is valid.

## How it works

The controller talks to the Dragonfly manager through its open API (`/oapi/v1/jobs`), authenticated
with a personal access token that has the `job` scope. Every operation on the P2P network is a manager
job that the controller creates and then polls until it reaches a terminal state:

| Operation  | Manager job   | When                                                                 |
| ---------- | ------------- | -------------------------------------------------------------------- |
| Distribute | `preheat`     | First reconcile, spec change, refresh due, retry after failure       |
| Verify     | `get_task`    | Verification interval elapsed                                        |
| Purge      | `delete_task` | TTL elapsed with `Purge`, source removed, deletion with `Purge`      |

A `Distribute` operation creates one preheat job per source. The tasks the preheat reports (the file
URLs, or the image layers) are recorded in the status together with their task IDs. Task IDs are
computed by the controller the same way the schedulers compute them, including the blob digest based
IDs used when `enableTaskIDBasedBlobDigest` is set, so that later `Verify` and `Purge` operations
address exactly the tasks that were distributed. `Verify` and `Purge` create one job per task, in
batches, so that a large image does not trip the manager's job rate limit.

Each source of a dataset moves through these states:

```text
Pending -> Distributing -> Ready <-> Degraded
                |            |
                v            v
              Failed      Expired / Purging -> Purged
```

- A failed distribution is retried with exponential backoff up to `distribution.retryLimit` times, then
  the source is `Failed` until the spec changes.
- Verification marks a source `Degraded` when a task is missing from every peer, and `Ready` again when
  it is found.
- When the TTL elapses the expire action runs. With `Purge` the data is removed from the peers and the
  source becomes `Expired`; with `Retain` the source is only marked `Expired`. A refresh that is due
  later re-distributes an expired source.
- Sources removed from the spec are purged first when the deletion policy is `Purge`, then dropped from
  the status.
- A finalizer keeps a deleted `Dataset` around until its data is purged when the deletion policy is
  `Purge`. A purge that fails for good is reported and does not block the deletion.

The dataset status summarises the sources in a `phase` (`Pending`, `Distributing`, `Ready`, `Degraded`,
`Failed`, `Expired`, `Purging` or `Suspended`) and in the conditions `Ready`, `Distributed`, `Verified`
and `Expired`. It also reports the last distribution and verification times, and the next refresh,
verification and expiration times. Events are recorded for every job created and every operation that
finishes.

Setting `spec.suspend` stops the controller from starting new operations while in-flight jobs keep
being tracked.

## Deployment

The manifests in `deploy/data-controller` install the CRDs, the RBAC, a `ConfigMap` with the controller
configuration, a `Secret` for the manager token and a two-replica `Deployment` with leader election.

1. Create a personal access token with the `job` scope in the manager console.
2. Put it in `deploy/data-controller/secret.yaml` and set the manager endpoint in
   `deploy/data-controller/configmap.yaml`.
3. Apply the manifests.

```shell
kubectl apply -k deploy/data-controller
kubectl apply -f deploy/data-controller/samples/dataset-files.yaml
kubectl get datasets -A
```

The controller image is built with `make docker-build-datacontroller` and the binary with
`make build-datacontroller`.

### Configuration

The configuration file defaults to `/etc/dragonfly/datacontroller.yaml`. Every key can be overridden by
an environment variable prefixed with `DATACONTROLLER_`, for example `DATACONTROLLER_MANAGER_TOKEN`.

```yaml
server:
  logLevel: info
manager:
  endpoint: http://dragonfly-manager.dragonfly-system.svc.cluster.local:8080
  tokenFile: /etc/dragonfly/manager/token
  requestTimeout: 30s
controller:
  namespace: ""
  leaderElection: true
  maxConcurrentReconciles: 4
  pollInterval: 15s
  resyncInterval: 10m
  retryBackoff: 30s
  maxRetryBackoff: 10m
metrics:
  enable: true
  addr: ":8000"
```

- `manager.endpoint` and `manager.token` (or `manager.tokenFile`) are the default manager. A dataset can
  override them with `spec.managerRef`.
- `controller.namespace` restricts the controller to one namespace; empty watches every namespace.
- `controller.pollInterval` is how often in-flight jobs are polled and `controller.resyncInterval` how
  often every dataset is reconciled without an event.
- `controller.retryBackoff` and `controller.maxRetryBackoff` bound the delay between retries of a
  failed operation.

### Metrics

The metrics endpoint exposes the standard controller-runtime metrics plus:

- `dragonfly_datacontroller_jobs_created_total{type,result}`: manager jobs created.
- `dragonfly_datacontroller_operations_completed_total{type,outcome}`: dataset operations finished.

## Development

The API types live in `datacontroller/apis/data/v1alpha1`. After changing them, regenerate the deepcopy
functions, the CRDs and the RBAC role:

```shell
make generate-datacontroller
```

The reconcilers live in `datacontroller/controllers` and are tested against a fake API server and a
fake manager:

```shell
go test ./datacontroller/...
```
