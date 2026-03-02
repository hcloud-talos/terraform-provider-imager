resource "imager_image" "example" {
  image_url     = "https://example.com/hcloud-amd64.raw.xz"
  architecture  = "x86"
  location      = "fsn1"
  server_type   = "cx23"
  description   = "Talos snapshot image"
  labels        = { os = "talos" }
  timeouts      = { create = "10m" }
}

