package orasclient

import (
	"context"
	"fmt"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"

	"github.com/uor-framework/client/registryclient"
)

type orasClient struct {
	registryOpts content.RegistryOptions
	copyOpts     []oras.CopyOpt
	fileStore    *content.File
	ref          string
}

var _ registryclient.Client = &orasClient{}

// GatherDescriptors loads files to create OCI descriptors.
func (c *orasClient) GatherDescriptors(files ...string) ([]ocispec.Descriptor, error) {
	fromFile := content.NewFile("")
	descs, err := loadFiles(fromFile, files...)
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
	if err := c.fileStore.Load(configDesc, config); err != nil {
		return configDesc, fmt.Errorf("unable to load new manifest config: %w", err)
	}
	return configDesc, nil
}

// GenerateManifest creates and stores a manifest.
// This is generated from the config descriptor and artifact descriptors.
func (c *orasClient) GenerateManifest(configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) error {
	manifest, manifestDesc, err := content.GenerateManifest(&configDesc, manifestAnnotations, descriptors...)
	if err != nil {
		return fmt.Errorf("unable to create manifest: %w", err)
	}

	return c.fileStore.StoreManifest(c.ref, manifestDesc, manifest)
}

// Execute performs the copy of OCI artifacts.
func (c *orasClient) Execute(ctx context.Context) error {
	to, err := content.NewRegistry(c.registryOpts)
	if err != nil {
		return fmt.Errorf("could not create registry target: %w", err)
	}
	desc, err := oras.Copy(ctx, c.fileStore, c.ref, to, "", c.copyOpts...)
	if err != nil {
		return err
	}
	fmt.Printf("Artifact published with digest %#v to %s\n", desc.Digest, c.ref)
	return nil
}

// Execute performs the copy of OCI artifacts.
func (c *orasClient) MapPaths(workDir string, descriptors ...ocispec.Descriptor) error {
	for _, desc := range descriptors {
		fpath := workDir + "/" + desc.Annotations["org.opencontainers.image.title"]
		c.fileStore.MapPath(desc.Annotations["org.opencontainers.image.title"], fpath)
	}
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
