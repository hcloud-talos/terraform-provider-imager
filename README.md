<div align="center">
  <br>
  <img src="https://github.com/hcloud-talos/terraform-hcloud-talos/blob/main/.idea/icon.png?raw=true" alt="Terraform - Hcloud - Talos" width="200"/>
  <h1>Terraform Provider: imager</h1>
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/hcloud-talos/terraform-provider-imager?logo=github">
  <p>Upload Talos disk images into Hetzner Cloud and turn them into reusable snapshot images.</p>
  <p>
    <a href="https://hetzner.cloud/?ref=9EF3RYocQW8y">New to Hetzner? Get 20EUR credit (and support this project)!</a>
  </p>
  <p>
    <a href="https://www.buymeacoffee.com/mrclrchtr"><img src="https://img.buymeacoffee.com/button-api/?text=Buy%20me%20a%20coffee&emoji=&slug=mrclrchtr&button_colour=FFDD00&font_colour=000000&font_family=Cookie&outline_colour=000000&coffee_colour=ffffff" alt="Buy me a coffee" /></a>
  </p>
  <p>If this provider saved you time or money, consider supporting ongoing maintenance.</p>
</div>

---

This repository contains a [Terraform provider](https://registry.terraform.io/providers/hcloud-talos/imager) for importing Talos `*.raw.xz` disk images into Hetzner Cloud.
It does this by creating a temporary rescue server, writing the image to disk, snapshotting it, and cleaning up the temporary resources afterwards.

- Use it when you want Talos snapshots managed from Terraform instead of building them manually.
- It fits especially well with Talos Image Factory URLs and the `siderolabs/talos` provider.
- The provider currently exposes a single resource: `imager_image`.

> [!WARNING]
> This provider is still alpha quality and not heavily battle-tested yet.
> Review created Hetzner Cloud resources and costs after apply and destroy.

---

## Goals

| Goal                          | Status | Description                                                                              |
|-------------------------------|--------|------------------------------------------------------------------------------------------|
| Terraform-native image import | OK     | Create Hetzner snapshot images directly from Terraform.                                  |
| Talos-focused workflow        | OK     | Designed around Talos `raw.xz` disk images, especially Image Factory outputs.            |
| Minimal provider surface      | OK     | Keeps the provider focused on image upload and lifecycle management.                     |
| Safe cleanup by default       | OK     | Temporary upload server resources are removed automatically unless debugging is enabled. |

## How It Works

1. Terraform reads a public `https://...raw.xz` URL.
2. The provider creates a temporary Hetzner rescue server.
3. The image is written to the root disk and converted into a snapshot.
4. The snapshot is kept, and temporary upload resources are deleted.

This is mainly useful as a building block for cluster modules such as [terraform-hcloud-talos](https://github.com/hcloud-talos/terraform-hcloud-talos), where snapshot creation needs to happen before cluster provisioning.

## Configuration

The provider uses `HCLOUD_TOKEN` by default:

```hcl
provider "imager" {}
```

Or configure the token explicitly:

```hcl
provider "imager" {
  token = var.hcloud_token
}
```

## Usage

### Minimal Example

```hcl
resource "imager_image" "example" {
  image_url    = "https://example.com/hcloud-amd64.raw.xz"
  architecture = "x86"

  location    = "fsn1"
  server_type = "cx23"

  description = "Talos snapshot image"

  labels = {
    os = "talos"
  }
}
```

### Talos Image Factory Workflow

A common setup is to generate the disk image URL with Terraform Provider Talos and then hand that URL to `imager_image`:

```hcl
data "talos_image_factory_urls" "hcloud_amd64" {
  talos_version = var.talos_version
  schematic_id  = talos_image_factory_schematic.this.id
  platform      = "hcloud"
  architecture  = "amd64"
}

resource "imager_image" "talos_x86" {
  image_url    = data.talos_image_factory_urls.hcloud_amd64.urls.disk_image
  architecture = "x86"

  labels = {
    os      = "talos"
    version = var.talos_version
  }
}
```

See [examples/talos-image-factory/main.tf](/Users/mrclrchtr/Development/mrclrchtr/terraform-provider-imager/examples/talos-image-factory/main.tf) for the complete example.

## Important Behavior

### Delete Semantics

Destroying `imager_image` deletes the resulting Hetzner snapshot image.

If that image is shared across multiple clusters or environments, protect it:

```hcl
lifecycle {
  prevent_destroy = true
}
```

### Input Expectations

- `image_url` must be a public HTTPS URL.
- The provider is intended for Talos `*.raw.xz` disk images.
- `architecture` must be `x86` or `arm`.
- `location` defaults to `fsn1` if not set.

### Debugging

Set `debug_skip_cleanup = true` only when you need to inspect failed upload runs manually. This intentionally leaves temporary server resources behind.

## Documentation

- Provider docs: [docs/index.md](/Users/mrclrchtr/Development/mrclrchtr/terraform-provider-imager/docs/index.md)
- Resource docs: [docs/resources/image.md](/Users/mrclrchtr/Development/mrclrchtr/terraform-provider-imager/docs/resources/image.md)

## Testing

### Unit Tests

```bash
mise run test
```

### Acceptance Tests

Acceptance tests create real Hetzner Cloud resources and are billable.

Required environment:

- `TF_ACC=1`
- `HCLOUD_TOKEN`
- `IMAGER_TEST_IMAGE_URL`

Run them with the pinned toolchain:

```bash
HCLOUD_TOKEN=... IMAGER_TEST_IMAGE_URL=... mise run test:acc
```

Or with 1Password Secrets in Environments:

```bash
op run --environment <environment-uuid> -- env TF_ACC=1 mise exec -- go test ./... -run TestAcc -count=1 -v
```

## Cleanup

If an upload fails at the wrong time, temporary servers or SSH keys can remain in your Hetzner project. The helper CLI can attempt cleanup of resources labeled by the upstream upload library:

```bash
HCLOUD_TOKEN=... mise run cleanup
```

You can also run the command directly:

```bash
HCLOUD_TOKEN=... go run ./cmd/imager-cleanup
```

## Development

Install the pinned toolchain and hooks:

```bash
mise install
mise run hooks
```

Useful commands:

```bash
mise run fix
mise run check
mise run build
mise run docs:gen
```

For local Terraform testing with `dev_overrides`, see [dev/local-terraform/README.md](/Users/mrclrchtr/Development/mrclrchtr/terraform-provider-imager/dev/local-terraform/README.md).

## Credits

This provider builds on top of the excellent [apricote/hcloud-upload-image](https://github.com/apricote/hcloud-upload-image) project by [@apricote](https://github.com/apricote).
Thanks for creating and maintaining it.
