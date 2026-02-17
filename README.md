# Terraform Provider: `imager`

Uploads Talos disk images (`*.raw.xz`) into Hetzner Cloud by creating a temporary rescue server, writing the image to the root disk, snapshotting it, and cleaning up.

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

## Cleanup

Failures can leave temporary servers or ssh keys behind. The upstream library labels temp resources with `apricote.de/created-by=hcloud-upload-image`.

You can attempt cleanup via:

```bash
HCLOUD_TOKEN=... go run ./cmd/imager-cleanup
```

## Credits

This provider builds on the excellent `hcloud-upload-image` project by @apricote. Thanks for creating and maintaining it.
