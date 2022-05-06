package registryclient

import (
	"context"
	"fmt"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

// Client defines methods to publish content as artifacts
type Client interface {
	GatherDescriptors(...string) ([]ocispec.Descriptor, error)
	GenerateConfig(map[string]string) (ocispec.Descriptor, error)
	GenerateManifest(ocispec.Descriptor, map[string]string, ...ocispec.Descriptor) error
	Execute(context.Context) error
}

type orasClient struct {
	registryOpts content.RegistryOptions
	copyOpts     []oras.CopyOpt
	fileStore    *content.File
	ref          string
}

var _ Client = &orasClient{}

func NewORASClient(ref string, copyOpts []oras.CopyOpt, registryOpts content.RegistryOptions) Client {
	return &orasClient{
		ref:          ref,
		copyOpts:     copyOpts,
		registryOpts: registryOpts,
	}
}

func (c *orasClient) GatherDescriptors(files ...string) ([]ocispec.Descriptor, error) {
	fromFile := content.NewFile("")
	descs, err := loadFiles(fromFile, files...)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}
	c.fileStore = fromFile
	return descs, nil
}

func (c *orasClient) GenerateConfig(configAnnotations map[string]string) (ocispec.Descriptor, error) {
	config, configDesc, err := content.GenerateConfig(configAnnotations)
	if err != nil {
		return configDesc, fmt.Errorf("unable to create new manifest config: %w", err)
	}
	if err := c.fileStore.Load(configDesc, config); err != nil {
		return configDesc, fmt.Errorf("unable to load new manifest config: %w", err)
	}
	return configDesc, nil
}

func (c *orasClient) GenerateManifest(configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) error {
	manifest, manifestDesc, err := content.GenerateManifest(&configDesc, manifestAnnotations, descriptors...)
	if err != nil {
		return fmt.Errorf("unable to create manifest: %w", err)
	}

	return c.fileStore.StoreManifest(c.ref, manifestDesc, manifest)
}

func (c *orasClient) Execute(ctx context.Context) error {
	to, err := content.NewRegistry(c.registryOpts)
	if err != nil {
		return fmt.Errorf("could not create registry target: %w", err)
	}
	desc, err := oras.Copy(ctx, c.fileStore, c.ref, to, "", c.copyOpts...)
	if err != nil {
		return err
	}
	fmt.Printf("Artifact published at %#v\n", desc.Digest)
	return nil
}

func loadFiles(store *content.File, files ...string) ([]ocispec.Descriptor, error) {
	var descs []ocispec.Descriptor
	for _, fileRef := range files {
		filename, mediaType := parseFileRef(fileRef, "")
		name := filepath.Clean(filename)
		if !filepath.IsAbs(name) {
			// convert to slash-separated path unless it is absolute path
			name = filepath.ToSlash(name)
		}
		desc, err := store.Add(name, mediaType, filename)
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}
	return descs, nil
}
