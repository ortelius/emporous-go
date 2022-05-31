package orasclient

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"

	"github.com/uor-framework/client/registryclient"
)

const uorMediaType = "application/vnd.uor.config.v1+json"

type orasClient struct {
	registryOpts content.RegistryOptions
	copyOpts     []oras.CopyOpt
	fileStore    *content.File
	ref          string
}

var _ registryclient.Client = &orasClient{}

// GatherDescriptors loads files to create OCI descriptors.
func (c *orasClient) GatherDescriptors(mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	fromFile := content.NewFile("")
	descs, err := loadFiles(fromFile, mediaType, files...)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}
	c.fileStore = fromFile
	return descs, nil
}

// GenerateConfig creates and stores a config.
// The config descriptor is returned for manifest generation.
func (c *orasClient) GenerateConfig(configAnnotations map[string]string) (ocispec.Descriptor, error) {
	config, configDesc, err := content.GenerateConfig(configAnnotations)
	if err != nil {
		return configDesc, fmt.Errorf("unable to create new manifest config: %w", err)
	}
	configDesc.MediaType = uorMediaType
	if err := c.fileStore.Load(configDesc, config); err != nil {
		return configDesc, fmt.Errorf("unable to load new manifest config: %w", err)
	}
	return configDesc, nil
}

// GenerateManifest creates and stores a manifest.
// This is generated from the config descriptor and artifact descriptors.
func (c *orasClient) GenerateManifest(configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) (ocispec.Descriptor, error) {
	manifest, manifestDesc, err := content.GenerateManifest(&configDesc, manifestAnnotations, descriptors...)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("unable to create manifest: %w", err)
	}

	if err := c.fileStore.StoreManifest(c.ref, manifestDesc, manifest); err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, err
}

// Execute performs the copy of OCI artifacts.
func (c *orasClient) Execute(ctx context.Context) (ocispec.Descriptor, error) {
	to, err := content.NewRegistry(c.registryOpts)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create registry target: %w", err)
	}
	desc, err := oras.Copy(ctx, c.fileStore, c.ref, to, "", c.copyOpts...)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	return desc, nil
}

func loadFiles(store *content.File, mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	var descs []ocispec.Descriptor
	var skipMediaTypeDetection bool
	var err error

	if mediaType != "" {
		skipMediaTypeDetection = true
	}
	for _, fileRef := range files {
		name := filepath.Clean(fileRef)
		if !filepath.IsAbs(name) {
			// convert to slash-separated path unless it is absolute path
			name = filepath.ToSlash(name)
		}

		if !skipMediaTypeDetection {
			mediaType, err = getDefaultMediaType(fileRef)
			if err != nil {
				return nil, fmt.Errorf("file %q: error dectecting media type: %v", name, err)
			}
		}

		desc, err := store.Add(name, mediaType, fileRef)
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}
	return descs, nil
}

func getDefaultMediaType(file string) (string, error) {
	mType, err := mimetype.DetectFile(file)
	if err != nil {
		return "", err
	}
	return mType.String(), nil
}
