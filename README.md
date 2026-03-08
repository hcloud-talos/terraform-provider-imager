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

## Quick Start

The provider uses `HCLOUD_TOKEN` by default. If you prefer, you can also set `token` explicitly in the provider block.

```hcl
terraform {
  required_version = ">= 1.9.0"

  required_providers {
    imager = {
      source  = "hcloud-talos/imager"
      version = "~> 0.1"
    }
  }
}

provider "imager" {}

resource "imager_image" "talos_x86" {
  image_url    = "https://factory.talos.dev/image/<schematic-id>/v1.12.4/hcloud-amd64.raw.xz"
  architecture = "x86"
}
```

If you omit optional arguments, the provider still creates a usable snapshot:

- `location` defaults to `fsn1`.
- `description` is generated automatically.
- Default snapshot labels are added for `os`, `creator`, and `arch`.

## How It Works

1. Terraform reads a public `https://...raw.xz` URL.
2. The provider creates a temporary Hetzner rescue server.
3. The image is written to the root disk and converted into a snapshot.
4. The snapshot is kept, and temporary upload resources are deleted.

This is mainly useful as a building block for cluster modules such as [terraform-hcloud-talos](https://github.com/hcloud-talos/terraform-hcloud-talos), where snapshot creation needs to happen before cluster provisioning.

## Talos Image Factory Workflow

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
    version = var.talos_version
  }
}
```

Architecture names differ between Talos and Hetzner:

- Talos `amd64` maps to provider `x86`.
- Talos `arm64` maps to provider `arm`.

See [examples/talos-image-factory/main.tf](examples/talos-image-factory/main.tf) for the complete example.

## Behavior and Defaults

### Input Requirements

- `image_url` must be a public `https://` URL.
- `image_url` must end in `.raw.xz`.
- The provider is intended for Talos `*.raw.xz` disk images, especially Hetzner Cloud images from Image Factory.
- `architecture` must be `x86` or `arm`.
- `labels` must follow Hetzner Cloud label rules. In particular, label values cannot contain `/`.

### Snapshot Metadata

- Default labels are always applied: `os=talos`, `creator=hcloud-talos-imager`, and `arch=<architecture>`.
- User-supplied `labels` are merged on top of those defaults.
- If `description` is omitted, the provider derives one from `labels.version` when present, otherwise from the Talos version embedded in the image URL when possible.
- The resource exposes `id`, `image_id`, and `effective_labels` after creation.

### Lifecycle Semantics

- Destroying `imager_image` deletes the resulting Hetzner snapshot image.
- Protect shared snapshots with `lifecycle { prevent_destroy = true }`.
- Changing `image_url`, `architecture`, `location`, `server_type`, or `debug_skip_cleanup` recreates the snapshot by running the upload again.
- Changing `description` updates the existing snapshot in place. `labels` are managed on the snapshot without a re-upload.
- Set `debug_skip_cleanup = true` only when you need to inspect a failed upload manually. This intentionally leaves temporary resources behind.
- Use the `timeouts` block if you need longer create, read, or delete timeouts for slow uploads.

```hcl
lifecycle {
  prevent_destroy = true
}

timeouts {
  create = "20m"
}
```

## Documentation

- Provider docs: [docs/index.md](docs/index.md)
- Resource docs: [docs/resources/image.md](docs/resources/image.md)

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

For provider development, install the pinned toolchain and hooks:

```bash
mise install
mise run hooks
```

Useful local commands:

```bash
mise run test
mise run test:acc
mise run fix
mise run check
mise run build
mise run docs:gen
```

For local Terraform testing with `dev_overrides`, see [dev/local-terraform/README.md](dev/local-terraform/README.md).

Acceptance tests create real Hetzner Cloud resources and are billable. They require `TF_ACC=1`, `HCLOUD_TOKEN`, and `IMAGER_TEST_IMAGE_URL`.

## Credits

This provider builds on top of the excellent [apricote/hcloud-upload-image](https://github.com/apricote/hcloud-upload-image) project by [@apricote](https://github.com/apricote).
Thanks for creating and maintaining it.
