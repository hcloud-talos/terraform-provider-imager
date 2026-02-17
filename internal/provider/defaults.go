package provider

import "maps"

func defaultLabels(architecture string) map[string]string {
	return map[string]string{
		"os":      "talos",
		"creator": "hcloud-talos/imager",
		"arch":    architecture,
	}
}

func mergeLabels(defaults, user map[string]string) map[string]string {
	out := maps.Clone(defaults)
	for k, v := range user {
		out[k] = v
	}
	return out
}
