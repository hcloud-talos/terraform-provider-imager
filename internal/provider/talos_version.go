package provider

import (
	"net/url"
	"regexp"
	"strings"
)

var talosVersionSegmentRegexp = regexp.MustCompile(`^v\d+\.\d+\.\d+([-.].+)?$`)

func extractTalosVersionFromImageURL(parsedURL *url.URL) string {
	if parsedURL == nil {
		return ""
	}

	segments := strings.Split(parsedURL.Path, "/")
	for _, segment := range segments {
		if talosVersionSegmentRegexp.MatchString(segment) {
			return segment
		}
	}

	return ""
}
