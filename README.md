# Terraform Provider: `imager`

Uploads Talos disk images (`*.raw.xz`) into Hetzner Cloud by creating a temporary rescue server, writing the image to the root disk, snapshotting it, and cleaning up.

> [!WARNING]
> This project is **alpha** quality and **not well-tested yet**. Use at your own risk, and double-check Hetzner Cloud resources and costs after apply/destroy.

## Configuration

The provider uses `HCLOUD_TOKEN` by default:

```hcl
provider "imager" {}
```

Or configure explicitly:

```hcl
provider "imager" {
  token = var.hcloud_token
}
```

## Resource: `imager_image`

```hcl
resource "imager_image" "talos_x86" {
  image_url    = data.talos_image_factory_urls.this.urls.disk_image
  architecture = "x86"

  location    = "fsn1"
  server_type = "cx23"

  labels = {
    os      = "talos"
    version = var.talos_version
  }

  timeouts {
    create = "10m"
  }
}
```

### Delete semantics warning

On `terraform destroy`, `imager_image` deletes the snapshot image from Hetzner Cloud. If you reuse images across multiple clusters, consider:

```hcl
lifecycle {
  prevent_destroy = true
}
```

## Using Talos Image Factory (Terraform Provider Talos)

This provider expects a public `https://...*.raw.xz` URL. A common workflow is to generate it with Terraform Provider Talos.

See `terraform-provider-imager/examples/with-talos/main.tf`.

## Testing

### Acceptance tests

Acceptance tests create real Hetzner Cloud resources (billable) and require:

- `TF_ACC=1`
- `HCLOUD_TOKEN`
- `IMAGER_TEST_IMAGE_URL` (public `https://...*.raw.xz` URL; e.g. `https://factory.talos.dev/image/376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba/v1.12.4/hcloud-amd64.raw.xz`)

Run them with environment variables:

```bash
TF_ACC=1 HCLOUD_TOKEN=... IMAGER_TEST_IMAGE_URL=... go test ./... -run TestAcc -count=1 -v
```

Or run them with 1Password CLI “Secrets in Environments” (environment `cq2r5uieu3ymytht2e3exxzuu4`):

```bash
op run --environment cq2r5uieu3ymytht2e3exxzuu4 -- env TF_ACC=1 go test ./... -run TestAcc -count=1 -v
```

## Cleanup

Failures can leave temporary servers or ssh keys behind. The upstream library labels temp resources with `apricote.de/created-by=hcloud-upload-image`.

You can attempt cleanup via:

```bash
HCLOUD_TOKEN=... go run ./cmd/imager-cleanup
```

## Development

Tooling is pinned via `mise.toml`:

```bash
mise install
mise run check
mise run docs:gen
```

### Local Terraform testing (dev overrides)

Install the provider binary into your `GOBIN` (or default Go bin dir):

```bash
go install ./cmd/terraform-provider-imager
```

Then configure Terraform CLI to use that directory (example `~/.terraformrc`):

```hcl
provider_installation {
  dev_overrides {
    "hcloud-talos/imager" = "/absolute/path/to/go/bin"
  }
  direct {}
}
```

## Releasing

This repo uses GoReleaser (`.goreleaser.yml`) and GitHub Actions (`.github/workflows/release.yml`) to publish Terraform Registry-compatible release artifacts.

- Tag a release like `v0.1.0`
- Configure repo secrets: `GPG_PRIVATE_KEY` and `PASSPHRASE`

## Credits

This provider builds on the excellent `hcloud-upload-image` project by @apricote. Thanks for creating and maintaining it.
