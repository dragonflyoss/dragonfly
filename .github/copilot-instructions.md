# Dragonfly P2P File Distribution System

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

Dragonfly is a P2P file distribution and image acceleration system written in Go 1.25.5. It consists of multiple components: manager (cluster management and web portal), scheduler (download optimization), dfget (P2P download client), dfcache (P2P cache operations), and dfstore (object storage with P2P cache).

## Working Effectively

### Bootstrap and Build Process

Install required dependencies and build the project:

```bash
# Install Go tools (required for linting and testing)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
go install github.com/onsi/ginkgo/v2/ginkgo@latest
export PATH=$PATH:$(go env GOPATH)/bin

# Build all components (takes ~1.5 minutes)
make build-manager-server build-scheduler build-dfget build-dfcache build-dfstore
# NEVER CANCEL: Build takes 2 minutes. Set timeout to 5+ minutes.

# Verify build success
ls -la bin/linux_amd64/
./bin/linux_amd64/dfget version
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

**CRITICAL**: The full `make build` target will fail because `build-manager-console` requires a Node.js frontend that is not included in this repository. Always use the individual build targets listed above to build only the Go components.

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
./bin/linux_amd64/dfget --help
./bin/linux_amd64/dfcache --help
./bin/linux_amd64/dfstore --help
./bin/linux_amd64/scheduler --help
./bin/linux_amd64/manager --help

# Test version commands (validates build was successful)
./bin/linux_amd64/dfget version
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

**Important**: You cannot run the full Dragonfly system in a sandboxed environment without proper network configuration, certificates, and storage backends. The components require complex setup with Redis, databases, and network connectivity between peers.

## Validation Requirements

### Pre-commit Validation

Always run these commands before committing changes:

```bash
# NEVER CANCEL: Full precheck takes 8 minutes. Set timeout to 15+ minutes.
make fmt vet
golangci-lint run --timeout=10m
make build-manager-server build-scheduler build-dfget build-dfcache build-dfstore

# Test that binaries still work
./bin/linux_amd64/dfget version
./bin/linux_amd64/scheduler version
./bin/linux_amd64/manager version
```

### Scenario Testing

After making changes, always validate:

1. **Build Success**: All Go components build without errors
2. **Binary Functionality**: Version commands execute successfully
3. **Help Commands**: All help text displays correctly
4. **Linting**: No linting errors introduced

## Repository Structure

### Key Directories

- `cmd/`: Main entry points for each component (dfget, scheduler, manager, dfcache, dfstore)
- `pkg/`: Shared libraries and utilities
- `scheduler/`: Scheduler service implementation
- `manager/`: Manager service implementation
- `client/`: Client-side P2P logic
- `test/`: Unit and E2E test suites
- `hack/`: Build scripts and utilities
- `api/`: API definitions and generated code

### Build Artifacts

- `bin/linux_amd64/`: Built binaries (created by build process)
- `coverage.txt`: Test coverage reports

### Important Files

- `Makefile`: All build, test, and lint targets
- `go.mod`: Go 1.2.5 dependencies
- `.golangci.yml`: Linting configuration
- `.markdownlint.yml`: Markdown linting rules

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

## Common Tasks

### Viewing Build Targets

```bash
make help
```

### Clean Build

```bash
make clean
rm -rf bin/
```

### Adding Dependencies

```bash
go mod tidy
go mod download
```

## Troubleshooting

### Build Issues

- **"build-manager-console" fails**: This is expected. The console frontend is not included. Use individual component build targets.
- **Missing tools**: Install golangci-lint and ginkgo as shown in bootstrap section.
- **Go version**: Requires Go 1.25.5 as specified in go.mod.

### Test Issues

- **Unit test failures**: Some tests may fail in sandboxed environments due to network/permission restrictions. This is expected.
- **E2E test failures**: Require Docker and complex setup. May not work in all environments.

### Runtime Issues

- **"command not found"**: Add `$(go env GOPATH)/bin` to PATH for installed Go tools.
- **Binary execution**: Built binaries are in `bin/linux_amd64/` directory.

## Development Workflow

1. **Make changes** to Go source files
2. **Format and vet**: `make fmt vet`
3. **Build**: Individual component build targets
4. **Test**: `./bin/linux_amd64/[component] version` to verify
5. **Lint**: `golangci-lint run --timeout=10m`
6. **Unit test**: `make test` (optional, may fail in sandbox)
7. **Commit** changes

Always build and validate binary functionality after code changes to ensure nothing is broken.
