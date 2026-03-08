#!/usr/bin/env bash
set -euo pipefail

bindir="$(go env GOBIN || true)"
if [[ -z "${bindir}" ]]; then
  bindir="$(go env GOPATH)/bin"
fi

provider_bin="${bindir}/terraform-provider-imager"
if [[ ! -x "${provider_bin}" ]]; then
  echo "Provider binary not found at: ${provider_bin}" >&2
  echo "Run: go install ./cmd/terraform-provider-imager" >&2
  exit 1
fi

cat > terraformrc <<EOF
provider_installation {
  dev_overrides {
    "hcloud-talos/imager" = "${bindir}"
  }
  direct {}
}
EOF

echo "Wrote ./terraformrc with dev override to: ${bindir}"

