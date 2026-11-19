# Scheduler Redis Cluster persistence

## Existing deployments

`database.redis.clusterKeyHashTag` defaults to `false` in both Manager and
Scheduler. The default preserves keys such as
`scheduler:scheduler-clusters:1:persistent-tasks:task-id`. Upgrading the binaries
alone does not switch the storage format. Standalone and Sentinel deployments
continue to use their existing keys and do not need metadata migration.

## Enable Cluster persistence

Redis Cluster support for the existing multi-key Lua scripts requires related
host, peer, task, and association keys to share a slot. Enable the new format
explicitly on every Manager and Scheduler using the same persistence Redis:

```yaml
database:
  redis:
    addrs:
      - redis-1.example:6379
      - redis-2.example:6379
    masterName: ""
    clusterKeyHashTag: true
    brokerDB: 0
    backendDB: 0
```

Manager also requires `database.redis.db: 0`. Multiple addresses with an empty
Sentinel master name select Cluster mode. Enabling `clusterKeyHashTag` for a
standalone or Sentinel connection fails configuration validation. A nonzero
database is rejected for Cluster connections.

With this option enabled, persistence keys use
`scheduler:scheduler-clusters:{1}:...`. Persistent task and peer enumeration,
including the Manager persistent-cache task API, scans the primary owning the
hash tag. The default legacy format does not provide multi-key Lua compatibility
with Redis Cluster.

## Upgrade and rollback

The format option is a storage setting, not a per-replica rollout flag. All
Manager and Scheduler replicas sharing persistence metadata must agree on its
value. The application deliberately does not dual-read or dual-write the two
formats: independently updating related records would break Lua atomicity.

For an existing deployment with legacy keys:

1. Upgrade the binaries with `clusterKeyHashTag: false`. Existing standalone and
   Sentinel deployments can stop here and retain normal rolling upgrades.
2. Before enabling the Cluster format, quiesce all persistence writers and
   maintenance workers and take a recoverable metadata backup. Plan migration of
   the complete host/peer/task relationship graph, including TTLs and values
   that reference keys; changing key names alone is not a complete migration.
3. Migrate offline, or explicitly rebuild disposable metadata under an agreed
   maintenance plan. Cross-slot `RENAME` is not a Redis Cluster migration method.
   Do not use `FLUSHDB` or `FLUSHALL` on a shared Redis deployment.
4. Enable the option consistently in Manager and Scheduler, verify task and
   peer enumeration, associations and persistent replica counts, then resume
   writers. Retain the backup until the cutover has been validated.

No automatic metadata migration is provided. Local file caches are not deleted
by the format option. After writes resume in the new format, changing the option
back would expose stale legacy metadata; rollback requires quiescing writers
and restoring or migrating a consistent metadata snapshot.

Deployments already using the `{schedulerClusterID}` format from an earlier
patch must explicitly set `clusterKeyHashTag: true` when adopting this revision.
The tagged key names remain the same. Leaving the option unset selects legacy
keys and is not a rollback procedure for those deployments.

## Sharding tradeoff

All related records for one scheduler cluster occupy one slot. This preserves
the atomicity of existing Lua scripts but does not distribute one scheduler
cluster's metadata across Redis primaries. Measure the owning primary's load
when sizing a large deployment; separate scheduler cluster IDs can use different
slots.
