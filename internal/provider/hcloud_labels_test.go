package provider

import "testing"

func TestValidateHcloudResourceLabels(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := validateHcloudResourceLabels(map[string]string{
			"os":      "talos",
			"creator": "hcloud-talos-imager",
			"arch":    "x86",
		})
		if err != nil {
			t.Fatalf("expected labels to be valid, got: %v", err)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		err := validateHcloudResourceLabels(map[string]string{
			"creator": "hcloud-talos/imager",
		})
		if err == nil {
			t.Fatal("expected invalid labels error, got nil")
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		err := validateHcloudResourceLabels(map[string]string{
			"-bad": "value",
		})
		if err == nil {
			t.Fatal("expected invalid labels error, got nil")
		}
	})
}
