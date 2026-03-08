# Local Terraform testing (dev override)

This directory lets you use the provider **from this repo** (a locally-built binary) without publishing to the Terraform Registry.

## 1) Build/install the provider binary

```bash
go install ../../cmd/terraform-provider-imager
```

## 2) Generate a Terraform CLI config with `dev_overrides`

```bash
./setup.sh
```

This writes `./terraformrc` pointing Terraform at your Go bin dir (where `terraform-provider-imager` was installed).
