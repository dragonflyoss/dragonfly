# Dragonfly P2P File Distribution System

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Project Overview

Dragonfly is a **CNCF Graduated** project that delivers efficient, stable, and secure data distribution and
acceleration powered by P2P technology. It provides an optional content-addressable filesystem that accelerates
OCI container launch. The project targets cloud-native architectures and large-scale delivery of files, container
images, OCI artifacts, AI/ML models, caches, logs, and dependencies.

- **Website**: <https://d7y.io>
- **Module path**: `d7y.io/dragonfly/v2`
- **Go version**: 1.25.5
- **Current version**: v2.4.0
- **License**: Apache 2.0

## Architecture

Dragonfly consists of the following components:

### Components in this repository

| Component     | Binary      | Default Ports       | Description                                                    |
| ------------- | ----------- | ------------------- | -------------------------------------------------------------- |
| **Manager**   | `manager`   | gRPC :65003, REST :8080 | Cluster management, peer lifecycle, dynamic configuration, web console |
| **Scheduler** | `scheduler` | gRPC :8002          | Optimizes P2P download scheduling and task routing             |

### Submodule components

| Submodule              | Repository                          | Description                                           |
| ---------------------- | ----------------------------------- | ----------------------------------------------------- |
| `client/`              | dragonflyoss/client                 | Rust-based dfdaemon (dfget, dfcache, dfstore, proxy)  |
| `manager/console/`     | dragonflyoss/console                | Node.js/React web console frontend for manager        |
| `deploy/helm-charts/`  | dragonflyoss/helm-charts            | Kubernetes Helm charts for deploying Dragonfly        |

### Data Flow

```
Client (dfdaemon/dfget) → Scheduler → Manager
                              ↓
                     Seed Peer / Origin Server
```

- **Manager**: stores cluster configuration in MySQL/PostgreSQL and caches in Redis; exposes both gRPC and REST APIs
- **Scheduler**: stateless gRPC service that manages the P2P DAG (task graph); stores transient state in Redis
- **Client (dfdaemon)**: runs on each node; intercepts HTTP/HTTPS requests and coordinates peer-to-peer downloads

## Repository Structure

### Key Directories

```
cmd/
  manager/          # Manager binary entry point (main.go)
  scheduler/        # Scheduler binary entry point (main.go)
  dependency/       # Shared CLI base options (base.Options) and plugin/version commands

manager/            # Manager service implementation
  config/           # Configuration structs and defaults (GRPC :65003, REST :8080)
  handlers/         # Gin HTTP request handlers
  router/           # HTTP router setup
  rpcserver/        # gRPC server implementation
  service/          # Business logic
  models/           # GORM database models
  database/         # Database initialization (MySQL, MariaDB, PostgreSQL)
  middlewares/      # Authentication, CORS, rate-limiting middleware
  job/              # Async job processing (dragonflyoss/machinery)
  metrics/          # Prometheus metrics
  searcher/         # Scheduler cluster selection logic
  auth/             # JWT authentication
  permission/       # Casbin RBAC authorization
  gc/               # Garbage collection

scheduler/          # Scheduler service implementation
  config/           # Configuration structs and defaults (port :8002)
  scheduling/       # Core P2P scheduling algorithms and evaluators
  resource/         # In-memory peer/task/host resource management
  rpcserver/        # gRPC server implementation
  service/          # Business logic
  announcer/        # Announces scheduler to manager
  job/              # Task job handling
  metrics/          # Prometheus metrics

pkg/                # Shared libraries
  auth/             # Authentication utilities
  balancer/         # gRPC load balancing
  cache/            # Cache utilities
  container/        # Container helpers
  dfnet/            # Network type definitions
  dfpath/           # Standard Dragonfly paths (/var/lib/dragonfly, /var/log/dragonfly)
  digest/           # Content digest (SHA-256, MD5, etc.)
  gc/               # Generic garbage collection framework
  graph/            # DAG implementation for task graphs
  idgen/            # ID generation for tasks, peers, hosts
  math/             # Math utilities
  net/              # Network utilities (IP, FQDN, HTTP)
  os/               # OS utilities
  redis/            # Redis client wrapper
  rpc/              # gRPC client wrappers (manager, scheduler, dfdaemon)
  slices/           # Generic slice utilities
  strings/          # String utilities
  structure/        # Data structure utilities
  time/             # Time utilities
  types/            # Shared type definitions

internal/           # Internal packages (not importable externally)
  dferrors/         # Dragonfly error types
  dflog/            # Structured logging (zap-based)
  dfplugin/         # Plugin system
  dynconfig/        # Dynamic configuration (polls manager for config updates)
  job/              # Shared job utilities
  ratelimiter/      # Rate limiter

api/
  manager/          # Swagger/OpenAPI generated docs for manager REST API

version/            # Version variables injected at build time

test/
  e2e/              # End-to-end tests using Ginkgo (require Docker)
  testdata/         # Shared test fixtures

deploy/
  docker-compose/   # Docker Compose setup (manager + scheduler + Redis + MySQL)
  helm-charts/      # Kubernetes Helm charts (submodule)

hack/               # Build and utility scripts
  build.sh          # Primary build script
  env.sh            # Build environment variables
  install.sh        # Binary installation script
  markdownlint.sh   # Markdown lint runner
```

### Build Artifacts

- `bin/linux_amd64/`: Built binaries (created by build process)
- `coverage.txt`: Accumulated test coverage reports

### Important Files

- `Makefile`: All build, test, and lint targets — run `make help` to see all
- `go.mod`: Go module definition (module `d7y.io/dragonfly/v2`, Go 1.25.5)
- `.golangci.yml`: Linting configuration (golangci-lint v2)
- `.markdownlint.yml`: Markdown linting rules

## Working Effectively

### Bootstrap and Build Process

Install required dependencies and build the project:

```bash
# Install Go tools (required for linting and testing)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
go install github.com/onsi/ginkgo/v2/ginkgo@latest
export PATH=$PATH:$(go env GOPATH)/bin

# Build Go server components (takes ~1.5 minutes)
make build-manager-server build-scheduler
# NEVER CANCEL: Build takes 2 minutes. Set timeout to 5+ minutes.

# Verify build success
ls -la bin/linux_amd64/
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

**CRITICAL**: The full `make build` target will fail because `build-manager-console` requires the
`manager/console` Node.js submodule which is not initialized by default. Always use
`build-manager-server` and `build-scheduler` to build only the Go server components.

### Testing

Run different test suites with appropriate timeouts:

```bash
# Format and vet code (takes ~20 seconds)
make fmt vet
# NEVER CANCEL: Set timeout to 2+ minutes.

# Run linting (takes ~1.5 minutes)
make markdownlint  # Takes ~5 seconds
golangci-lint run --timeout=10m  # Takes ~1.5 minutes
# NEVER CANCEL: Linting takes 2 minutes total. Set timeout to 5+ minutes.

# Run unit tests (takes ~3.5 minutes, may have some failures in fresh environment)
make test
# NEVER CANCEL: Unit tests take 4 minutes. Set timeout to 10+ minutes.
# Note: Some tests may fail in sandboxed environments, this is expected.

# Run E2E tests (requires Docker and longer setup)
make e2e-test
# NEVER CANCEL: E2E tests take 10+ minutes. Set timeout to 30+ minutes.
```

### Running Applications

Test the built applications:

```bash
# Test CLI help (validates binaries work correctly)
./bin/linux_amd64/scheduler --help
./bin/linux_amd64/manager --help

# Test version commands (validates build was successful)
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

**Important**: You cannot run the full Dragonfly system in a sandboxed environment without proper network
configuration, certificates, Redis, and a relational database. The components require complex setup with
network connectivity between peers.

## Validation Requirements

### Pre-commit Validation

Always run these commands before committing changes:

```bash
# NEVER CANCEL: Full precheck takes 8 minutes. Set timeout to 15+ minutes.
make fmt vet
golangci-lint run --timeout=10m
make build-manager-server build-scheduler

# Test that binaries still work
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

### Scenario Testing

After making changes, always validate:

1. **Build Success**: All Go components build without errors
2. **Binary Functionality**: Version commands execute successfully
3. **Help Commands**: All help text displays correctly
4. **Linting**: No linting errors introduced

## Timing Expectations

| Command             | Duration    | Timeout Recommendation |
| ------------------- | ----------- | ---------------------- |
| `make fmt vet`      | 20 seconds  | 2 minutes              |
| `make markdownlint` | 5 seconds   | 1 minute               |
| `golangci-lint run` | 1.5 minutes | 5 minutes              |
| Build (Go only)     | 1.5 minutes | 5 minutes              |
| `make test`         | 3.5 minutes | 10 minutes             |
| `make e2e-test`     | 10+ minutes | 30 minutes             |
| Full precheck       | 8 minutes   | 15 minutes             |

**CRITICAL**: NEVER CANCEL long-running commands. Builds and tests are CPU-intensive and require time to complete.

## Code Conventions

### License Header

Every Go source file must begin with the Apache 2.0 license header:

```go
/*
 *     Copyright <YEAR> The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
```

### Import Ordering

Imports must be organized in exactly four groups (enforced by `gci` linter):

```go
import (
    // 1. Standard library
    "context"
    "fmt"

    // 2. Third-party dependencies
    "github.com/gin-gonic/gin"

    // 3. d7y.io/api packages
    "d7y.io/api/v2/pkg/apis/scheduler/v1"

    // 4. d7y.io/dragonfly/v2 packages (this module)
    "d7y.io/dragonfly/v2/pkg/idgen"
)
```

### Error Handling

- Use `d7y.io/dragonfly/v2/internal/dferrors` for Dragonfly-specific error types
- Use `fmt.Errorf("context: %w", err)` for wrapping errors
- The `errcheck` linter enforces that all errors are handled

### Logging

- Use `d7y.io/dragonfly/v2/internal/dflog` (zap-based structured logger)
- Log fields use `zap.String`, `zap.Int`, etc.

### Configuration

- Configuration structs use both `yaml` and `mapstructure` struct tags
- All configs embed `base.Options` from `cmd/dependency/base` for shared CLI options
- Config files default to `/etc/dragonfly/manager.yaml` and `/etc/dragonfly/scheduler.yaml`
- Dynamic configuration is polled from the manager via gRPC

### gRPC

- Proto definitions are in the separate `d7y.io/api/v2` module
- gRPC client wrappers are in `pkg/rpc/` (manager, scheduler, dfdaemon clients)
- All gRPC services implement health checking (`grpc_health_probe`)

### REST API (Manager only)

- Built with Gin framework (`github.com/gin-gonic/gin`)
- Swagger docs generated with `swag` — run `make swag` to regenerate
- Generated files are in `api/manager/`
- Routes defined in `manager/router/`
- Handlers in `manager/handlers/`
- JWT authentication via `github.com/appleboy/gin-jwt/v2`
- RBAC authorization via Casbin (`github.com/casbin/casbin/v2`)

### Testing Patterns

- Unit tests use the standard `testing` package and live alongside source in `_test.go` files
- Mocks are generated with mockery and placed in `mocks/` subdirectories
- E2E tests use Ginkgo v2 and are located in `test/e2e/`
- Use `-race` flag for race condition detection (already set in `make test`)
- Test fixtures go in `testdata/` directories next to the tests

## Configuration Overview

### Manager Configuration (`/etc/dragonfly/manager.yaml`)

Key sections:

```yaml
server:
  grpc:
    port: { start: 65003, end: 65003 }  # gRPC port range
  rest:
    addr: ":8080"                         # REST API address
database:
  type: mysql                             # mysql | mariadb | postgres
  mysql:
    host: localhost
    port: 3306
    db: manager
cache:
  redis:
    addrs: ["localhost:6379"]
```

### Scheduler Configuration (`/etc/dragonfly/scheduler.yaml`)

Key sections:

```yaml
server:
  port: 8002                   # gRPC port
  advertisePort: 8002
scheduler:
  algorithm: default           # Scheduling algorithm
  backToSourceCount: 200       # Max peers that can download from origin
database:
  redis:
    addrs: ["localhost:6379"]
manager:
  addr: "localhost:65003"      # Manager gRPC address
```

## Deployment

### Docker Compose (Development)

```bash
cd deploy/docker-compose
# Edit config files in ./config/ as needed
docker compose up -d
```

Services started: Redis, MySQL, Manager (gRPC :65003, REST :8080), Scheduler (:8002), Client (seed + regular)

### Kubernetes (Helm)

The `deploy/helm-charts/` submodule contains Helm charts. Initialize it first:

```bash
git submodule update --init deploy/helm-charts
helm install dragonfly deploy/helm-charts/charts/dragonfly
```

### External Dependencies

| Dependency     | Purpose                               | Required by       |
| -------------- | ------------------------------------- | ----------------- |
| Redis          | Cache, job queue, scheduler state     | Manager, Scheduler |
| MySQL/MariaDB  | Persistent storage (config, metadata) | Manager           |
| PostgreSQL     | Alternative to MySQL                  | Manager           |

## Common Tasks

### View All Build Targets

```bash
make help
```

### Clean Build

```bash
make clean
rm -rf bin/
```

### Add/Update Dependencies

```bash
go mod tidy
go mod download
```

### Regenerate Swagger API Docs

```bash
make swag
# Output: api/manager/
```

### Initialize Submodules

```bash
git submodule update --init --recursive
```

## Troubleshooting

### Build Issues

- **`build-manager-console` fails**: Expected — the `manager/console` frontend submodule is not initialized.
  Use `build-manager-server` and `build-scheduler` instead.
- **Missing tools**: Install golangci-lint and ginkgo as shown in the bootstrap section.
- **Go version**: Requires Go 1.25.5 as specified in `go.mod`.

### Test Issues

- **Unit test failures**: Some tests may fail in sandboxed environments due to network/permission
  restrictions. This is expected.
- **E2E test failures**: Require Docker and a running Dragonfly cluster. May not work in all environments.

### Runtime Issues

- **`command not found`**: Add `$(go env GOPATH)/bin` to `PATH` for installed Go tools.
- **Binary execution**: Built binaries are in `bin/linux_amd64/` directory.

## Development Workflow

1. **Make changes** to Go source files
2. **Format and vet**: `make fmt vet`
3. **Build**: `make build-manager-server build-scheduler`
4. **Verify binaries**: `./bin/linux_amd64/manager version` and `./bin/linux_amd64/scheduler version`
5. **Lint**: `golangci-lint run --timeout=10m`
6. **Unit test**: `make test` (optional, may fail in sandbox)
7. **Commit** changes

Always build and validate binary functionality after code changes to ensure nothing is broken.
