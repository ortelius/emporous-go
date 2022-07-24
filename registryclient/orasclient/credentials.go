package orasclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/credentials"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// AuthStore contains authentication configuration
// file information for registry interactions.
type AuthStore struct {
	configs []*configfile.ConfigFile
}

// NewAuthStore returns a new authentication store instance
// with credential information collected.
func NewAuthStore(configPaths ...string) (*AuthStore, error) {
	if len(configPaths) == 0 {
		// No config path passed, attempt to load default configuration
		// from well-known locations.
		cfg, err := loadDefaultConfig()
		if err != nil {
			return nil, err
		}

		return &AuthStore{
			configs: []*configfile.ConfigFile{cfg},
		}, nil
	}

	var configs []*configfile.ConfigFile
	for _, path := range configPaths {
		cfg, err := loadConfigFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		configs = append(configs, cfg)
	}

	return &AuthStore{
		configs: configs,
	}, nil
}

// Credential iterates all the config files, returns the first non-empty
// credential in a best-effort way.
func (s *AuthStore) Credential(_ context.Context, registry string) (auth.Credential, error) {
	for _, c := range s.configs {
		authConf, err := c.GetCredentialsStore(registry).Get(registry)
		if err != nil {
			return auth.EmptyCredential, err
		}
		cred := auth.Credential{
			Username:     authConf.Username,
			Password:     authConf.Password,
			AccessToken:  authConf.RegistryToken,
			RefreshToken: authConf.IdentityToken,
		}
		if cred != auth.EmptyCredential {
			return cred, nil
		}
	}
	return auth.EmptyCredential, nil
}

// loadConfigFile reads the credential-related configuration
// from the given path.
func loadConfigFile(path string) (*configfile.ConfigFile, error) {
	cfg := configfile.New(path)
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return cfg, err
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := cfg.LoadFromReader(file); err != nil {
		return nil, err
	}

	if !cfg.ContainsAuth() {
		cfg.CredentialsStore = credentials.DetectDefaultStore(cfg.CredentialsStore)
	}

	return cfg, nil
}

// loadDefaultConfig attempts to load credentials in the
// default Docker location, then the default Podman location.
func loadDefaultConfig() (*configfile.ConfigFile, error) {
	dir := config.Dir()
	dockerConfigJSON := filepath.Join(dir, config.ConfigFileName)
	cfg := configfile.New(dockerConfigJSON)

	switch _, err := os.Stat(dockerConfigJSON); {
	case err == nil:
		cfg, err = config.Load(dir)
		if err != nil {
			return cfg, err
		}
	case os.IsNotExist(err):
		podmanConfig := filepath.Join(xdg.RuntimeDir, "containers/auth.json")
		cfg, err = loadConfigFile(podmanConfig)
		if err != nil {
			return cfg, err
		}
	}

	if !cfg.ContainsAuth() {
		cfg.CredentialsStore = credentials.DetectDefaultStore(cfg.CredentialsStore)
	}

	return cfg, nil
}
