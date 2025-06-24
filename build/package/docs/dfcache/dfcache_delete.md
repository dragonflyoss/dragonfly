% DFCACHE(1) Version v2.2.0 | Frivolous "Dfcache" Documentation

# NAME

**dfcache delete** — delete file from P2P cache system

# SYNOPSIS

Delete file from P2P cache system.

```shell
dfcache delete <-i cid> [flags]
```

## OPTIONS

```shell
  -i, --cid string            content or cache ID, e.g. sha256 digest of the content
      --config string         the path of configuration file with yaml extension name, default is /etc/dragonfly/dfcache.yaml, it can also be set by env var: DFCACHE_CONFIG
      --console               whether logger output records to the stdout
      --logdir string         Dfcache log directory
  -t, --tag string            different tags for the same cid will be recognized as different  files in P2P network
      --timeout duration      Timeout for this cache operation, 0 is infinite
      --workhome string       Dfcache working directory
  -h, --help   help for delete
```

# SEE ALSO

- [dfcache](dfcache.md) - the P2P cache client of dragonfly
