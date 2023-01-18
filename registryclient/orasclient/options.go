package orasclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry/remote/auth"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/content/memory"

	"github.com/emporous/emporous-go/content"
	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/registryclient"
)

// ClientOption is a function that configures
// options on the client config.
type ClientOption func(o *ClientConfig) error

// ClientConfig contains configuration data for the registry client.
type ClientConfig struct {
	configs    []string
	credFn     func(context.Context, string) (auth.Credential, error)
	prePullFn  func(context.Context, string) error
	plainHTTP  bool
	insecure   bool
	cache      content.Store
	copyOpts   oras.CopyOptions
	attributes model.Matcher
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

	// Setup auth client based on config inputs
	authClient := &auth.Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.insecure,
				},
			},
		},
		Cache: auth.NewCache(),
	}

	if config.credFn != nil {
		authClient.Credential = config.credFn
	} else {
		store, err := NewAuthStore(config.configs...)
		if err != nil {
			return nil, err
		}
		authClient.Credential = store.Credential
	}

	client.authClient = authClient
	client.plainHTTP = config.plainHTTP
	client.copyOpts = config.copyOpts
	client.destroy = destroy
	client.cache = config.cache
	client.attributes = config.attributes
	client.prePullFn = config.prePullFn

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
