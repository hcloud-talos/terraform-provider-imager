# Contributing

## Prerequisites

- `mise` (installs pinned Go, `hk`, `golangci-lint`)

## Local checks

```bash
mise install
mise run check
mise run docs:check
```

## Tests

Unit tests:

```bash
go test ./...
```

Acceptance tests (billable Hetzner Cloud resources):

```bash
TF_ACC=1 HCLOUD_TOKEN=... IMAGER_TEST_IMAGE_URL=... go test ./... -run TestAcc -count=1 -v
```

## Security

Please do not open public issues for security reports. See `SECURITY.md`.

