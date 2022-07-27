package cli

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
)

func pushBlob(ctx context.Context, mediaType string, blob []byte, target oras.Target) (ocispec.Descriptor, error) {
	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    digest.FromBytes(blob),
		Size:      int64(len(blob)),
	}
	return desc, target.Push(ctx, desc, bytes.NewReader(blob))
}

func generateManifest(configDesc ocispec.Descriptor, layers ...ocispec.Descriptor) ([]byte, error) {
	manifest := ocispec.Manifest{
		Config:    configDesc,
		Layers:    layers,
		Versioned: specs.Versioned{SchemaVersion: 2},
	}
	return json.Marshal(manifest)
}
