# Repository Guidelines

## Project Structure & Module Organization

- `cmd/terraform-provider-imager/`: Terraform provider entrypoint (serves `registry.terraform.io/hcloud-talos/imager`).
- `cmd/imager-cleanup/`: CLI helper to remove leaked temporary Hetzner resources.
- `internal/provider/`: provider implementation (`imager_image`) plus unit and acceptance tests.
- `examples/with-talos/`: example Terraform config using Talos Image Factory URLs.
- `README.md`: usage, configuration, and operational notes.

## Build, Test, and Development Commands

- `mise install`: installs pinned tools from `mise.toml` (Go, `hk`, `golangci-lint`).
- `hk install`: installs Git hooks defined in `hk.pkl` (format + lint on commit).
- `hk check` / `hk fix`: runs the repo checks locally (Go format, `go mod tidy`, `golangci-lint`).
- `go test ./...`: runs unit tests.
- `TF_ACC=1 HCLOUD_TOKEN=... IMAGER_TEST_IMAGE_URL=... go test ./... -run TestAcc -count=1`: runs acceptance tests (creates real Hetzner Cloud resources).
- `go build ./cmd/terraform-provider-imager`: builds the provider binary.
- `HCLOUD_TOKEN=... go run ./cmd/imager-cleanup`: attempts cleanup of temp servers/SSH keys left behind after failures.

## Coding Style & Naming Conventions

- Prefer `hk fix` before pushing (runs `gofmt`, `go mod tidy`, and `golangci-lint --fix`).
- Keep packages small and focused; prefer clear types over clever abstractions.
- Tests: `*_test.go` for unit tests; `*_acceptance_test.go` and `TestAcc*` for acceptance tests.

Tooling notes:
- Hook runner: `hk` (config: `hk.pkl`). References: [Getting Started](https://hk.jdx.dev/getting_started.html), [Builtins: golangci-lint](https://hk.jdx.dev/builtins.html#golangci-lint).

## Testing Guidelines

- Unit tests should be deterministic and avoid network calls.
- Acceptance tests must guard on env vars (see `TF_ACC`, `HCLOUD_TOKEN`, `IMAGER_TEST_IMAGE_URL`) and should clean up resources they create.

## Commit & Pull Request Guidelines

- Commit messages follow Conventional Commits style (e.g., `feat: ...`, `fix: ...`, `docs: ...`) as in Git history.
- PRs should include: what changed, why, how to test (commands + env vars), and any cost/impact notes (acceptance tests create billable resources).

## Security & Configuration Tips

- Never commit secrets. Use environment variables (not hardcoded tokens) for Hetzner credentials.
- Prefer running acceptance tests in a dedicated Hetzner project to avoid accidental deletion of shared snapshots.
