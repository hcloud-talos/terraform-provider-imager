terraform {
  required_providers {
    imager = {
      source = "hcloud-talos/imager"
    }
  }
}

provider "imager" {
  # token = var.hcloud_token
}

