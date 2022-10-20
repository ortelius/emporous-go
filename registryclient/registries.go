package registryclient

import (
	"regexp"
	"strings"
)

// This configuration is slightly modified and paired down version of the registries.conf.
// Source https://github.com/containers/image/blob/main/pkg/sysregistriesv2/system_registries_v2.go.
// More information on why this does not just use the `containers/system_registries_v2` library.
// While this library has a lot of overlapping functionality, it has more functionality than we
// need, and it makes sense to use the `containers` registry client which we are not. Search registries
// will eventually be a used in this library, but will be resolved and related to collection attributes
// and not short names.

// Endpoint describes a remote location of a registry.
type Endpoint struct {
	// The endpoint's remote location.
	Location string `mapstructure:"location" json:"location"`
	// If true, certs verification will be skipped.
	SkipTLS bool `mapstructure:"skipTLS" json:"skipTLS"`
	// If true, the client will use HTTP to
	// connect to the registry.
	PlainHTTP bool `mapstructure:"plainHTTP" json:"plainHTTP"`
}

// Registry represents a registry.
type Registry struct {
	// Prefix is used for endpoint matching.
	Prefix string `mapstructure:"prefix" json:"prefix"`
	// A registry is an Endpoint too
	Endpoint `mapstructure:",squash" json:",inline"`
}

// RegistryConfig is a configuration to configure multiple
// registry endpoints.
type RegistryConfig struct {
	Registries []Registry `mapstructure:"registries" json:"registries"`
}

// FindRegistry returns the registry from the registry config that
// matches the reference.
func FindRegistry(registryConfig RegistryConfig, reference string) (*Registry, error) {
	reg := Registry{}
	prefixLen := 0

	for _, r := range registryConfig.Registries {
		match := r.Prefix
		if match == "" {
			match = r.Location
		}
		prefixExp, err := regexp.Compile(validPrefix(match))
		if err != nil {
			return nil, err
		}
		if prefixExp.MatchString(reference) {
			if len(match) > prefixLen {
				reg = r
				prefixLen = len(match)
			}
		}
	}
	if prefixLen != 0 {
		return &reg, nil
	}
	return nil, nil
}

// validPrefix will check the registry prefix value
// and return a valid regex.
func validPrefix(regPrefix string) string {
	if strings.HasPrefix(regPrefix, "*") {
		return strings.Replace(regPrefix, "*", ".*", -1)
	}
	return regPrefix
}
