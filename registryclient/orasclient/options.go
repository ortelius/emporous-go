package orasclient

import (
	"context"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"oras.land/oras-go/v2"
	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/content/memory"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/registryclient"
)

// ClientOption is a function that configures
// options on the client config.
type ClientOption func(o *ClientConfig) error

// ClientConfig contains configuration data for the registry client.
type ClientConfig struct {
	outputDir string
	configs   []string
	plainHTTP bool
	insecure  bool
	cache     content.Store
	copyOpts  oras.CopyOptions
}

func (c *ClientConfig) apply(options []ClientOption) error {
	for _, option := range options {
		if err := option(c); err != nil {
			return err
		}
	}
	return nil
}

// NewClient returns a new ORAS client implementation
func NewClient(options ...ClientOption) (registryclient.Client, error) {
	client := &orasClient{}

	config := &ClientConfig{}
	config.copyOpts = oras.DefaultCopyOptions
	if err := config.apply(options); err != nil {
		return client, err
	}

	var once sync.Once
	destroy := func() (destroyErr error) {
		once.Do(func() {
			destroyErr = client.artifactStore.Close()
		})

		return
	}

	client.insecure = config.insecure
	client.configs = config.configs
	client.plainHTTP = config.plainHTTP
	client.copyOpts = config.copyOpts
	client.outputDir = config.outputDir
	client.destroy = destroy
	client.cache = config.cache

	// We are not allowing this to be configurable since
	// oras file stores turn artifacts into descriptors in
	// specific way we want to reuse.
	client.artifactStore = file.NewWithFallbackStorage("", memory.New())

	return client, nil
}

// WithAuthConfigs adds configuration files
// with registry authorization information.
func WithAuthConfigs(configs []string) ClientOption {
	return func(config *ClientConfig) error {
		config.configs = configs
		return nil
	}
}

// SkipTLSVerify disables TLS certificate checking.
func SkipTLSVerify(insecure bool) ClientOption {
	return func(config *ClientConfig) error {
		config.insecure = insecure
		return nil
	}
}

// WithPlainHTTP uses the HTTP protocol with the registry.
func WithPlainHTTP(plainHTTP bool) ClientOption {
	return func(config *ClientConfig) error {
		config.plainHTTP = plainHTTP
		return nil
	}
}

// WithCache uses the provided storage a cache to be used
// with remote resources. It is the responsibility of the caller
// to perform any clean up actions.
func WithCache(store content.Store) ClientOption {
	return func(config *ClientConfig) error {
		config.cache = store
		return nil
	}
}

// WithPostCopy applies a function to a descriptor after copying it.
// This sets the oras.CopyOptions.PostCopy function.
func WithPostCopy(postFn func(ctx context.Context, desc ocispec.Descriptor) error) ClientOption {
	return func(config *ClientConfig) error {
		config.copyOpts.PostCopy = postFn
		return nil
	}
}

// WithPreCopy applies a function to a descriptor before copying it.
// This sets the oras.CopyOptions.PreCopy function.
func WithPreCopy(preFn func(ctx context.Context, desc ocispec.Descriptor) error) ClientOption {
	return func(config *ClientConfig) error {
		config.copyOpts.PreCopy = preFn
		return nil
	}
}

// WithSuccessorFn adds a function to find the child node of the current node if exists.
// This sets the oras.CopyOptions.FindSuccessor function.
func WithSuccessorFn(successorFn func(ctx context.Context, fetcher orascontent.Fetcher, desc ocispec.Descriptor) ([]ocispec.Descriptor, error)) ClientOption {
	return func(config *ClientConfig) error {
		config.copyOpts.FindSuccessors = successorFn
		return nil
	}
}
