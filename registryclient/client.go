package registryclient

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Client defines methods to interact with OCI artifacts
// in various contexts.
type Client interface {
	// GatherDescriptors loads files to create OCI descriptors with a specific
	// media type.
	GatherDescriptors(context.Context, string, ...string) ([]ocispec.Descriptor, error)
	// GenerateConfig creates and stores a config.
	// The config descriptor is returned for manifest generation.
	GenerateConfig(context.Context, []byte, map[string]string) (ocispec.Descriptor, error)
	// GenerateManifest creates and stores a manifest for an image reference.
	// This is generated from the config descriptor and artifact descriptors.
	GenerateManifest(context.Context, string, ocispec.Descriptor, map[string]string, ...ocispec.Descriptor) (ocispec.Descriptor, error)
	// Execute performs the copy of OCI artifacts.
	// The image reference and copy action are expected as inputs.
	Execute(context.Context, string, ActionType) (ocispec.Descriptor, error)
	// Destroy cleans up any on-disk resources used to track descriptors.
	Destroy() error
}

// ActionType defines what actions (e.g. push, push, etc...) the execute method should perform.
type ActionType int

const (
	// TypeInvalid is the default action.
	// It is invalid because the action must
	// be explicitly set.
	TypeInvalid ActionType = iota
	// TypePush action pushes from a local location
	// to a remote location.
	TypePush
	// TypePull action pulls from a remote location to
	// a local location.
	TypePull
)
