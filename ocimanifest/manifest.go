package ocimanifest

import (
	"encoding/json"
	"io"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// FetchSchemaLinks fetches schema information from a given input.
func FetchSchemaLinks(input io.Reader) (string, []string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
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

// ResolveCollectionLinks finds linked collection references from a given input.
func ResolveCollectionLinks(input io.Reader) ([]string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return nil, err
	}
	links, ok := manifest.Annotations[AnnotationCollectionLinks]
	if !ok || len(links) == 0 {
		return nil, ErrNoCollectionLinks
	}
	return strings.Split(links, Separator), nil
}
