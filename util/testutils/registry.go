package testutils

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"oras.land/oras-go/v2/content"
)

// NewRegistry returns a handler which implements a mock registry with a v3
// attributes endpoint.
func NewRegistry(t *testing.T, blobs [][]byte, manifests [][]byte) http.Handler {
	blobsByDigest := map[string][]byte{}
	for _, blob := range blobs {
		d := digest.FromBytes(blob)
		blobsByDigest[d.String()] = blob
	}

	manifestsByDigest := map[string][]byte{}
	var manifestDescs []ocispec.Descriptor
	for _, manifest := range manifests {
		d := digest.FromBytes(manifest)
		manifestsByDigest[d.String()] = manifest
		manifestDesc := content.NewDescriptorFromBytes(ocispec.MediaTypeArtifactManifest, manifest)
		manifestDescs = append(manifestDescs, manifestDesc)
	}

	testIndex := ocispec.Index{
		MediaType: ocispec.MediaTypeImageIndex,
		Manifests: manifestDescs,
	}
	testIndexJSON, err := json.Marshal(testIndex)
	require.NoError(t, err)
	d := digest.FromBytes(testIndexJSON)
	registryFN := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected access: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		switch {
		case strings.HasPrefix(r.URL.Path, "/v2/test/blobs/"):
			parts := strings.Split(r.URL.Path, "/")
			digest := parts[len(parts)-1]
			content, ok := blobsByDigest[digest]
			if !ok {
				t.Errorf("failed to find blob")
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Docker-Content-Digest", digest)
			if _, err := w.Write(content); err != nil {
				t.Errorf("failed to write %q: %v", r.URL, err)
			}
		case strings.HasPrefix(r.URL.Path, "/v2/test/manifests/"):
			if accept := r.Header.Get("Accept"); !strings.Contains(accept, ocispec.MediaTypeImageManifest) {
				t.Errorf("manifest not convertable: %s", accept)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			parts := strings.Split(r.URL.Path, "/")
			digest := parts[len(parts)-1]
			content, ok := manifestsByDigest[digest]
			if !ok {
				t.Errorf("failed to find blob")
			}
			w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
			w.Header().Set("Docker-Content-Digest", digest)
			if _, err := w.Write(content); err != nil {
				t.Errorf("failed to write %q: %v", r.URL, err)
			}
		case strings.HasPrefix(r.URL.Path, "/v2/attributes"):
			values := r.URL.Query()
			if !values.Has("attributes") && !values.Has("digests") && !values.Has("links") {
				t.Errorf("wrong query type")
			}
			w.Header().Set("Content-Type", ocispec.MediaTypeImageIndex)
			w.Header().Set("Docker-Content-Digest", d.String())
			if _, err := w.Write(testIndexJSON); err != nil {
				t.Errorf("failed to write %q: %v", r.URL, err)
			}
		default:
			t.Errorf("unexpected access: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	return registryFN
}
