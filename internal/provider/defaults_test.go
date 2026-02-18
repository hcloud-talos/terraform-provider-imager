package provider

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestMergeLabels_UserOverridesDefaults(t *testing.T) {
	defaults := map[string]string{
		"os":      "talos",
		"creator": "hcloud-talos/imager",
	}
	user := map[string]string{
		"creator": "override",
		"team":    "platform",
	}

	got := mergeLabels(defaults, user)
	if got["creator"] != "override" {
		t.Fatalf("expected creator override, got %q", got["creator"])
	}
	if got["team"] != "platform" {
		t.Fatalf("expected team label, got %q", got["team"])
	}
	if got["os"] != "talos" {
		t.Fatalf("expected os label, got %q", got["os"])
	}
}

func TestDefaultLabels_AreValidForHetznerResources(t *testing.T) {
	rawLabels := map[string]interface{}{}
	for key, value := range defaultLabels("x86") {
		rawLabels[key] = value
	}

	ok, err := hcloud.ValidateResourceLabels(rawLabels)
	if err != nil {
		t.Fatalf("expected valid default labels, got error: %v", err)
	}
	if !ok {
		t.Fatal("expected valid default labels")
	}
}
