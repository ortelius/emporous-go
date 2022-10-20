package orasclient

import (
	"context"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry/remote/auth"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/content/memory"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/registryclient"
)

// ClientOption is a function that configures
// options on the client config.
type ClientOption func(o *ClientConfig) error

// ClientConfig contains configuration data for the registry client.
type ClientConfig struct {
	outputDir      string
	configs        []string
	credFn         func(context.Context, string) (auth.Credential, error)
	plainHTTP      bool
	skipTLSVerify  bool
	cache          content.Store
	copyOpts       oras.CopyOptions
	attributes     model.Matcher
	registryConfig registryclient.RegistryConfig
	prePullFn  func(context.Context, string) error
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

	client.authCache = auth.NewCache()
	client.plainHTTP = config.plainHTTP
	client.skipTLSVerify = config.skipTLSVerify
	client.copyOpts = config.copyOpts
	client.destroy = destroy
	client.cache = config.cache
	client.attributes = config.attributes
	client.registryConf = config.registryConfig
	client.prePullFn = config.prePullFn

	if config.credFn != nil {
		client.credFn = config.credFn
	} else {
		store, err := NewAuthStore(config.configs...)
		if err != nil {
			return nil, err
		}
		client.credFn = store.Credential
	}


	// We are not allowing this to be configurable since
	// oras file stores turn artifacts into descriptors in
	// specific way we want to reuse.
	client.artifactStore = file.NewWithFallbackStorage("", memory.New())

	return client, nil
}

// WithCredentialFunc overrides the default credential function. Using this option will override
// WithAuthConfigs.
func WithCredentialFunc(credFn func(context.Context, string) (auth.Credential, error)) ClientOption {
	return func(config *ClientConfig) error {
		config.credFn = credFn
		return nil
	}
}

// WithAuthConfigs adds configuration files
// with registry authorization information.
func WithAuthConfigs(configs []string) ClientOption {
	return func(config *ClientConfig) error {
		config.configs = configs
		return nil
	}
}

// WithRegistryConfig defines the configuration for specific registry
// endpoints. If specified, the configuration for a found registry
// will override WithSkipTLSVerify and WithPlainHTTP.
func WithRegistryConfig(registryConf registryclient.RegistryConfig) ClientOption {
	return func(config *ClientConfig) error {
		config.registryConfig = registryConf
		return nil
	}
}

// SkipTLSVerify disables TLS certificate checking.
func SkipTLSVerify(skipTLSVerify bool) ClientOption {
	return func(config *ClientConfig) error {
		config.skipTLSVerify = skipTLSVerify
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

// WithPrePullFunc applies a function to a reference before pulling it to a content
// store.
func WithPrePullFunc(prePullFn func(context.Context, string) error) ClientOption {
	return func(config *ClientConfig) error {
		config.prePullFn = prePullFn
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

// WithPullableAttributes adds a filter when pulling blobs that allows non-matching
// blobs to be skipped.
func WithPullableAttributes(filter model.Matcher) ClientOption {
	return func(config *ClientConfig) error {
		config.attributes = filter
		return nil
	}
}
