package provider

import "testing"

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
