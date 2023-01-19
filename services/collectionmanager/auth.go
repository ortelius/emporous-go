package collectionmanager

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/auth"

	managerapi "github.com/emporous/emporous-go/api/services/collectionmanager/v1alpha1"
)

// authConfig wraps the AuthConfig from the manager API.
type authConfig struct {
	auth *managerapi.AuthConfig
}

// Credential returns the credential specified from the AuthConfig if the host matches.
func (s *authConfig) Credential(_ context.Context, registry string) (auth.Credential, error) {
	if s.auth == nil {
		return auth.EmptyCredential, nil
	}
	if s.auth.RegistryHost != "" {
		if registry != s.auth.RegistryHost {
			return auth.EmptyCredential, nil
		}
	}

	cred := auth.Credential{
		Username:     s.auth.Username,
		Password:     s.auth.Password,
		AccessToken:  s.auth.AccessToken,
		RefreshToken: s.auth.RefreshToken,
	}
	if cred != auth.EmptyCredential {
		return cred, nil
	}

	return auth.EmptyCredential, nil
}
