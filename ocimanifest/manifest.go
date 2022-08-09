package ocimanifest

import (
	"context"
	"encoding/json"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/registryclient"
)

// FetchSchema fetches schema information for a reference.
func FetchSchema(ctx context.Context, reference string, client registryclient.Remote) (string, []string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return "", nil, err
	}

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manBytes).Decode(&manifest); err != nil {
		return "", nil, err
	}

	schema, ok := manifest.Annotations[AnnotationSchema]
	if !ok {
		return "", nil, ErrNoKnownSchema
	}

	links, ok := manifest.Annotations[AnnotationSchemaLinks]
	if !ok {
		return schema, nil, nil
	}

	return schema, []string{links}, nil
}
