package provider

import "testing"

func TestValidateAndParseImageURL(t *testing.T) {
	_, err := validateAndParseImageURL("https://example.com/hcloud-amd64.raw.xz")
	if err != nil {
		t.Fatalf("expected valid url, got %v", err)
	}

	_, err = validateAndParseImageURL("http://example.com/hcloud-amd64.raw.xz")
	if err == nil {
		t.Fatalf("expected error for non-https scheme")
	}

	_, err = validateAndParseImageURL("https://example.com/hcloud-amd64.iso")
	if err == nil {
		t.Fatalf("expected error for non-.raw.xz path")
	}
}

func TestMapArchitecture(t *testing.T) {
	_, label, err := mapArchitecture("x86")
	if err != nil || label != "x86" {
		t.Fatalf("expected x86, got label=%q err=%v", label, err)
	}
	_, label, err = mapArchitecture("arm")
	if err != nil || label != "arm" {
		t.Fatalf("expected arm, got label=%q err=%v", label, err)
	}
	_, _, err = mapArchitecture("amd64")
	if err == nil {
		t.Fatalf("expected error for invalid architecture")
	}
}
