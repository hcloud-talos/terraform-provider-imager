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

## Releases

Releases are cut automatically after successful pushes to `main`. If a `next` branch is used, it produces prereleases.

Versioning is derived from Conventional Commits:

- `feat`: minor release
- `fix`, `perf`, `docs`, `revert`: patch release
- `BREAKING CHANGE` or `type!`: major release

The semantic-release workflow creates the Git tag, and that tag triggers GoReleaser to publish the Terraform provider artifacts.

## Security

Please do not open public issues for security reports. See `SECURITY.md`.
