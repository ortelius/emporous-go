package orasclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/uor-framework/client/registryclient"
)

const uorMediaType = "application/vnd.uor.config.v1+json"

type orasClient struct {
	insecure  bool
	plainHTTP bool
	configs   []string
	copyOpts  oras.CopyOptions
	fileStore *file.Store
	destroy   func() error
	outputDir string
}

var _ registryclient.Client = &orasClient{}

// GatherDescriptors loads files to create OCI descriptors.
func (c *orasClient) GatherDescriptors(ctx context.Context, mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	c.init()
	descs, err := loadFiles(ctx, c.fileStore, mediaType, files...)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}
	return descs, nil
}

// GenerateConfig creates and stores a config.
// The config descriptor is returned for manifest generation.
func (c *orasClient) GenerateConfig(ctx context.Context, config []byte, configAnnotations map[string]string) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}

	configDesc := ocispec.Descriptor{
		MediaType:   uorMediaType,
		Digest:      digest.FromBytes(config),
		Size:        int64(len(config)),
		Annotations: configAnnotations,
	}

	return configDesc, c.fileStore.Push(ctx, configDesc, bytes.NewReader(config))
}

// GenerateManifest creates and stores a manifest.
// This is generated from the config descriptor and artifact descriptors.
func (c *orasClient) GenerateManifest(ctx context.Context, ref string, configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}
	if descriptors == nil {
		descriptors = []ocispec.Descriptor{}
	}

	// Keep descriptor order deterministic
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Digest < descriptors[j].Digest
	})

	var packOpts oras.PackOptions
	packOpts.ConfigDescriptor = &configDesc
	packOpts.ManifestAnnotations = manifestAnnotations

	manifestDesc, err := oras.Pack(ctx, c.fileStore, descriptors, packOpts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, c.fileStore.Tag(ctx, manifestDesc, ref)
}

// Execute performs the copy of OCI artifacts.
func (c *orasClient) Execute(ctx context.Context, ref string, typ registryclient.ActionType) (ocispec.Descriptor, error) {
	var to, from oras.Target
	repo, err := c.setupRepo(ref)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create registry target: %w", err)
	}

	switch typ {
	case registryclient.TypePush:
		to = repo
		from = c.fileStore
	case registryclient.TypePull:
		c.fileStore = file.New(c.outputDir)
		to = c.fileStore
		from = repo
	case registryclient.TypeInvalid:
		return ocispec.Descriptor{}, errors.New("action type must be set")
	default:
		return ocispec.Descriptor{}, errors.New("unsupported action type")
	}

	desc, err := oras.Copy(ctx, from, ref, to, "", c.copyOpts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return desc, nil
}

// Destroy cleans up any on-disk resources used to track descriptors.
func (c *orasClient) Destroy() error {
	return c.destroy()
}

// init will initialize the file store
// if not set to avoid panics.
func (c *orasClient) init() {
	if c.fileStore == nil {
		c.fileStore = file.New("")
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

// setupRepo configures the client to access the remote repository.
func (c *orasClient) setupRepo(ref string) (*remote.Repository, error) {
	repo, err := remote.NewRepository(ref)
	if err != nil {
		return nil, fmt.Errorf("could not create registry target: %w", err)
	}
	repo.PlainHTTP = c.plainHTTP
	authC, err := c.authClient()
	if err != nil {
		return nil, err
	}
	repo.Client = authC
	return repo, nil
}

// authClient gather authorization information
// for registry access from provided and default configuration
// files.
func (c *orasClient) authClient() (*auth.Client, error) {
	client := &auth.Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.insecure,
				},
			},
		},
		Cache: auth.NewCache(),
	}

	store, err := NewAuthStore(c.configs...)
	if err != nil {
		return nil, err
	}
	client.Credential = store.Credential
	return client, nil
}

// loadFiles stores files in a file store and creates descriptors representing each file in the store.
func loadFiles(ctx context.Context, store *file.Store, mediaType string, files ...string) ([]ocispec.Descriptor, error) {
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

		desc, err := store.Add(ctx, name, mediaType, fileRef)
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}
	return descs, nil
}

// getDefaultMediaType detects the media type of the
// file based on content.
func getDefaultMediaType(file string) (string, error) {
	mType, err := mimetype.DetectFile(file)
	if err != nil {
		return "", err
	}
	return mType.String(), nil
}
