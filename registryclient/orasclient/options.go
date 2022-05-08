package orasclient

import (
	"github.com/uor-framework/client/registryclient"
)

// TODO(jpower432): Allow configuration for relevant ORAS copy options

type ClientOption func(o *ClientConfig) error

type ClientConfig struct {
	configs   []string
	plainHTTP bool
	insecure  bool
}

func (c *ClientConfig) apply(options []ClientOption) error {
	for _, option := range options {
		if err := option(c); err != nil {
			return err
		}
	}
	return nil
}

func NewClient(ref string, options ...ClientOption) (registryclient.Client, error) {
	client := &orasClient{
		ref: ref,
	}

	config := &ClientConfig{}
	if err := config.apply(options); err != nil {
		return client, err
	}

	client.registryOpts.Insecure = config.insecure
	client.registryOpts.Configs = config.configs
	client.registryOpts.PlainHTTP = config.plainHTTP
	return client, nil
}

func WithAuthConfigs(configs []string) ClientOption {
	return func(config *ClientConfig) error {
		config.configs = configs
		return nil
	}
}

func SkipTLSVerify(insecure bool) ClientOption {
	return func(config *ClientConfig) error {
		config.insecure = insecure
		return nil
	}
}

func WithPlainHTTP(plainHTTP bool) ClientOption {
	return func(config *ClientConfig) error {
		config.plainHTTP = plainHTTP
		return nil
	}
}
