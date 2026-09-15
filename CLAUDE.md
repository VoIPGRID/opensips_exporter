# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

A Prometheus exporter for OpenSIPS metrics. It talks to OpenSIPS 3.0+ over `mi_http`
(JSON-RPC over HTTP). The `mi_datagram` protocol (OpenSIPS ≤2.4) was dropped in 2.1.0.

## Commands

```bash
# Build
go build github.com/VoIPGRID/opensips_exporter

# Test
go test ./...

# Run a single test
go test ./opensips/... -run TestName

# Lint / code quality
make format       # gofmt
make vet          # go vet
make staticcheck  # staticcheck (requires: go install honnef.co/go/tools/cmd/staticcheck@latest)
make all          # format + vet + staticcheck + build

# Docker-based integration testing with a real OpenSIPS instance
make docker-test-build   # Build OpenSIPS test image
make docker-test-start   # Start OpenSIPS on port 8888
```

## Architecture

### Processor pattern

Each OpenSIPS module (core, shmem, dialog, etc.) has a dedicated processor in `processors/`. Processors implement `prometheus.Collector` (`Describe` + `Collect`) and register themselves in `init()` via `OpensipsProcessors` map in `processor.go`.

Key registration convention: each processor registers both individual statistic keys (e.g. `"core:rcv_requests"`) and a module wildcard key (e.g. `"core:"`). The HTTP handler uses the wildcard when `collect[]` query params are not specified.

### Adding a new processor

1. Create `processors/<module>_processor.go`
2. Define a `map[string]metric` for the module's statistics
3. Implement the processor struct with `Describe`/`Collect`
4. In `init()`, register each metric key and the module wildcard in `OpensipsProcessors`

### Protocol packages

- `opensips/opensips.go` — just the shared `Statistic` type
- `opensips/jsonrpc/opensips_jsonrpc.go` — `mi_http` JSON-RPC implementation, exposes
  `GetStatistics(targets ...string) (map[string]Statistic, error)`

### Testing

There are currently no unit tests. Use the Docker targets below to test against a real OpenSIPS
instance. Tests use the standard library `testing` package (no testify).
