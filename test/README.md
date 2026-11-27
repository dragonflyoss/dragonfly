# Test

Dragonfly includes unit tests and E2E tests.

## Unit tests

Unit tests is in the project directory.

### Running unit tests

```bash
make test
```

### Running unit tests with coverage reports

```bash
make test-coverage
```

## E2E tests

E2E tests is in `test/e2e` path.

### Running E2E tests

```bash
make e2e-test
```

### Seed cache recovery after Scheduler restart

Use the deployed E2E kind/Helm environment from `.github/workflows/e2e.yml`.
The test uses the existing Scheduler StatefulSet, all seed peers, an ordinary
client, and the file server. It does not require Manager and is included in
lightweight mode. The kube context must allow Pod reads, exec and port-forward,
and restarting the Scheduler StatefulSet.
`make e2e-test` runs the suite; it does not provision this environment.

Run only the cache-recovery regression:

```bash
go test -v -race -timeout 30m ./test/e2e -run '^TestE2E$' \
  -args -ginkgo.label-filter='scheduler && restart && seed-cache'
```

The test preloads a unique 32 MiB file on every seed, restarts every Scheduler,
waits for host registration, and downloads the file once from a cold client.
It verifies the file checksum, P2P download bytes, and zero additional source
piece bytes on the client and every seed. Metadata probes and failed source
requests are not counted by these byte counters. Peer restarts and counter
resets fail the test instead of being accepted as zero source traffic.

### Running E2E tests with coverage reports

```bash
make e2e-test-coverage
```

### Clean E2E tests environment

```bash
make clean-e2e-test
```
