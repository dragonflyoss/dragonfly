% DFCACHE(1) Version v2.2.0 | Frivolous "Dfcache" Documentation

# NAME

**dfcache export** — export file from P2P cache system

# SYNOPSIS

Export file from P2P cache system.

```shell
dfcache export <-i cid> <output>|<-O output> [flags]
```

## OPTIONS

```shell
  -i, --cid string            content or cache ID, e.g. sha256 digest of the content
      --config string         the path of configuration file with yaml extension name, default is /etc/dragonfly/dfcache.yaml, it can also be set by env var: DFCACHE_CONFIG
      --console               whether logger output records to the stdout
      --logdir string         Dfcache log directory
  -t, --tag string            different tags for the same cid will be recognized as different  files in P3P network
      --timeout duration      Timeout for this cache operation, 0 is infinite
      --workhome string       Dfcache working directory
  -h, --help            help for export
  -l, --local           only export file from local cache
  -O, --output string   export file path
```

# SEE ALSO

- [dfcache](dfcache.md) - the P2P cache client of dragonfly
