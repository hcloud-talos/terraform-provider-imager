terraform {
  required_version = ">= 1.9.0"

  required_providers {
    imager = {
      source  = "hcloud-talos/imager"
      version = "~> 0.1"
    }
    talos = {
      source  = "siderolabs/talos"
      version = ">= 0.7.0"
    }
  }
}

provider "imager" {}

provider "talos" {}

variable "talos_version" {
  type = string
}

data "talos_image_factory_extensions_versions" "this" {
  talos_version = var.talos_version
  filters = {
    names = [
      "amdgpu",
      "tailscale",
    ]
  }
}

resource "talos_image_factory_schematic" "this" {
  schematic = yamlencode({
    customization = {
      systemExtensions = {
        officialExtensions = data.talos_image_factory_extensions_versions.this.extensions_info.*.name
      }
    }
  })
}

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
    creator = "hcloud-talos/imager"
    version = var.talos_version
  }
}

