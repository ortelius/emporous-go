package options

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/uor-framework/uor-client-go/registryclient"
)

// Remote describes remote configuration options that can be set.
type Remote struct {
	SkipTLSVerify  bool
	PlainHTTP      bool
	RegistryConfig registryclient.RegistryConfig
}

// BindFlags binds options from a flag set to Remote options.
func (o *Remote) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.SkipTLSVerify, "skip-tls-verify", o.SkipTLSVerify, "disable TLS certificate verification when contacting registries")
	fs.BoolVar(&o.PlainHTTP, "plain-http", o.PlainHTTP, "use plain http and not https when contacting registries")
}

// LoadRegistryConfig loads the registry config from disk.
func (o *Remote) LoadRegistryConfig() error {
	viper.SetConfigName("registry-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.uor")
	err := viper.ReadInConfig()
	if err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFound) {
			return nil
		}
		return err
	}

	return viper.Unmarshal(&o.RegistryConfig)
}

// RemoteAuth describes remote authentication configuration options that can be set.
type RemoteAuth struct {
	Configs []string
}

// BindFlags binds options from a flag set to RemoteAuth options.
func (o *RemoteAuth) BindFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&o.Configs, "configs", "c", o.Configs, "auth config paths when contacting registries")
}
