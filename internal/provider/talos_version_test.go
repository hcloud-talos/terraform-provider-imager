package provider

import (
	"net/url"
	"testing"
)

func TestExtractTalosVersionFromImageURL(t *testing.T) {
	u, err := url.Parse("https://factory.talos.dev/image/abc/v1.12.4/hcloud-amd64.raw.xz")
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	if got := extractTalosVersionFromImageURL(u); got != "v1.12.4" {
		t.Fatalf("expected v1.12.4, got %q", got)
	}

	u, err = url.Parse("https://example.com/hcloud-amd64.raw.xz")
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	if got := extractTalosVersionFromImageURL(u); got != "" {
		t.Fatalf("expected empty version, got %q", got)
	}
}
