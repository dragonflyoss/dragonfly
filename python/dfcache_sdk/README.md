# Dragonfly Dfcache Python SDK

This is an initial, lightweight Python interface to Dragonfly's `dfcache` operations aimed at AI and data workflows.

## Features
- Import a local file into Dragonfly P2P cache
- Export a cached file by content ID (CID) to a local path
- Stat (existence check) for a CID (optionally local-only)
- Delete a cached CID
- Health check convenience

## Design
MVP implementation shells out to the existing `dfcache` CLI binary rather than using gRPC. This avoids proto generation overhead. A future iteration can switch to direct gRPC calls and piece streaming for advanced scenarios.

## Install
```bash
pip install -e python/dfcache_sdk
```
Ensure the Dragonfly `dfcache` binary is built:
```bash
make build-dfcache
```

## Usage
```python
from dragonfly_dfcache import DfCacheClient, NotFoundError

client = DfCacheClient()

cid = "sha256:abcdef1234567890"  # Example digest
source_file = "/data/model.bin"
export_path = "/tmp/model.bin"

# Import
client.import_cache(cid=cid, path=source_file)

# Stat
exists = client.stat(cid)
print("Exists?", exists)

# Export
client.export(cid=cid, output=export_path)

# Delete
client.delete(cid)
```

## Health
```python
client.check_health()  # True if dfcache CLI and dfdaemon responsive
```

## Roadmap
- Replace CLI subprocess calls with gRPC client (requires Python stubs)
- Add streaming progress callbacks for export/import
- Optional async API using asyncio subprocess
- Support rate limiting and advanced flags

## License
Apache 2.0
