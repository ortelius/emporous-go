package registryclient

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Client defines methods to publish content as artifacts
type Client interface {
	// GatherDescriptors loads files to create OCI descriptors with a specific
	// media type.
	GatherDescriptors(string, ...string) ([]ocispec.Descriptor, error)
	// GenerateConfig creates and stores a config.
	// The config descriptor is returned for manifest generation.
	GenerateConfig(map[string]string) (ocispec.Descriptor, error)
	// GenerateManifest creates and stores a manifest.
	// This is generated from the config descriptor and artifact descriptors.
	GenerateManifest(ocispec.Descriptor, map[string]string, ...ocispec.Descriptor) error
	// Execute performs the copy of OCI artifacts.
	Execute(context.Context) error
}
