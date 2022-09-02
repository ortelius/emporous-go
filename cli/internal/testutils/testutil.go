package testutils

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"

	"github.com/uor-framework/uor-client-go/content/layout"
)

func PushBlob(ctx context.Context, mediaType string, blob []byte, target oras.Target) (ocispec.Descriptor, error) {
	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    digest.FromBytes(blob),
		Size:      int64(len(blob)),
	}
	return desc, target.Push(ctx, desc, bytes.NewReader(blob))
}

func GenerateManifest(configDesc ocispec.Descriptor, annotations map[string]string, layers ...ocispec.Descriptor) ([]byte, error) {
	manifest := ocispec.Manifest{
		Config:      configDesc,
		Layers:      layers,
		Versioned:   specs.Versioned{SchemaVersion: 2},
		Annotations: annotations,
	}
	return json.Marshal(manifest)
}

// PrepCache will push a hello.txt artifact into the
// registry for retrieval. Uses methods from oras-go.
func PrepCache(ref string, cacheDir string, fileAnnotations map[string]string) error {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	ctx := context.TODO()

	ociStore, err := layout.NewWithContext(ctx, cacheDir)
	if err != nil {
		return err
	}
	layerDesc, err := PushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, ociStore)
	if err != nil {
		return err
	}
	if layerDesc.Annotations == nil {
		layerDesc.Annotations = map[string]string{}
	}
	for k, v := range fileAnnotations {
		layerDesc.Annotations[k] = v
	}
	layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

	config := []byte("{}")
	configDesc, err := PushBlob(ctx, ocispec.MediaTypeImageConfig, config, ociStore)
	if err != nil {
		return err
	}

	manifest, err := GenerateManifest(configDesc, nil, layerDesc)
	if err != nil {
		return err
	}

	manifestDesc, err := PushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, ociStore)
	if err != nil {
		return err
	}

	return ociStore.Tag(ctx, manifestDesc, ref)
}
