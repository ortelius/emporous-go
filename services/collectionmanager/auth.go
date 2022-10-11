package collectionmanager

import (
	"context"
	"fmt"
	"net/url"

	"oras.land/oras-go/v2/registry/remote/auth"

	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
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
	if s.auth.ServerAddress != "" {
		// Do not return the auth info when server address doesn't match.
		u, err := url.Parse(s.auth.ServerAddress)
		if err != nil {
			return auth.EmptyCredential, fmt.Errorf("parse server address: %w", err)
		}
		if registry != u.Host {
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
