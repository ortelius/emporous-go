package orasclient

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"oras.land/oras-go/pkg/target"

	"github.com/uor-framework/client/registryclient"
)

const uorMediaType = "application/vnd.uor.config.v1+json"

type orasClient struct {
	registryOpts content.RegistryOptions
	copyOpts     []oras.CopyOpt
	fileStore    *content.File
	outputDir    string
}

var _ registryclient.Client = &orasClient{}

// GatherDescriptors loads files to create OCI descriptors.
func (c *orasClient) GatherDescriptors(mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	c.init()
	descs, err := loadFiles(c.fileStore, mediaType, files...)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}
	return descs, nil
}

// GenerateConfig creates and stores a config.
// The config descriptor is returned for manifest generation.
func (c *orasClient) GenerateConfig(configAnnotations map[string]string) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}

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
func (c *orasClient) GenerateManifest(ref string, configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}

	manifest, manifestDesc, err := content.GenerateManifest(&configDesc, manifestAnnotations, descriptors...)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("unable to create manifest: %w", err)
	}

	if err := c.fileStore.StoreManifest(ref, manifestDesc, manifest); err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, err
}

// Execute performs the copy of OCI artifacts.
func (c *orasClient) Execute(ctx context.Context, ref string, typ registryclient.ActionType) (ocispec.Descriptor, error) {
	var to, from target.Target
	reg, err := content.NewRegistry(c.registryOpts)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create registry target: %w", err)
	}

	switch typ {
	case registryclient.TypePush:
		if err := c.checkFileStore(); err != nil {
			return ocispec.Descriptor{}, err
		}
		to = reg
		from = c.fileStore
	case registryclient.TypePull:
		c.fileStore = content.NewFile(c.outputDir)
		to = c.fileStore
		from = reg
	case registryclient.TypeInvalid:
		return ocispec.Descriptor{}, errors.New("action type must be set")
	default:
		return ocispec.Descriptor{}, errors.New("unsupported action type")
	}

	desc, err := oras.Copy(ctx, from, ref, to, "", c.copyOpts...)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	return desc, c.fileStore.Close()
}

// init will initialize the file store
// if not set to avoid panics.
func (c *orasClient) init() {
	if c.fileStore == nil {
		c.fileStore = content.NewFile("")
	}
}

// checkFileStore ensure that the file store
// has been initialized.
func (c *orasClient) checkFileStore() error {
	if c.fileStore == nil {
		return errors.New("file store uninitialized")
	}
	return nil
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
