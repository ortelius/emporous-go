//go:build !windows
// +build !windows

package registryclient

import "strings"

// QUESTION(jpower432): should we be determining media type a different way and
// passing to ORAS? Config file perhaps?
func parseFileRef(ref string, mediaType string) (string, string) {
	i := strings.LastIndex(ref, ":")
	if i < 0 {
		return ref, mediaType
	}
	return ref[:i], ref[i+1:]
}
