package queries

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"oras.land/oras-go/v2/registry/remote/auth"
)

type QueryParamsFn func(values url.Values)

// ResolveQuery send a query to the v3 registry.
func ResolveQuery(ctx context.Context, registryHost string, paramsFn QueryParamsFn, client *auth.Client, plainHTTP bool) ([]byte, error) {
	attributeURL, err := constructAttributesURL(registryHost, plainHTTP)
	if err != nil {
		return nil, err
	}

	queryParams := attributeURL.Query()
	paramsFn(queryParams)
	attributeURL.RawQuery = queryParams.Encode()

	ctx = auth.AppendScopes(ctx, auth.ActionPull)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, attributeURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// probe server range request ability.
	// Docker spec allows range header form of "Range: bytes=<start>-<end>".
	// The form of "Range: bytes=<start>-" is also acceptable.
	// However, the remote server may still not RFC 7233 compliant.
	// Reference: https://docs.docker.com/registry/spec/api/#blob
	req.Header.Set("Range", "bytes=0-")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func constructAttributesURL(baseURL string, plainHTTP bool) (*url.URL, error) {
	attributes := buildRegistryAttributesURL(baseURL, plainHTTP)
	return url.Parse(attributes)
}

// buildRegistryAttributesURL builds the URL for accessing the attributes API.
// Format: <scheme>://<registry>/v2/_attributes
func buildRegistryAttributesURL(baseURL string, plainHTTP bool) string {
	return fmt.Sprintf("%s://%s/v2/attributes", buildScheme(plainHTTP), baseURL)
}

// buildScheme returns HTTP scheme used to access the remote registry.
func buildScheme(plainHTTP bool) string {
	if plainHTTP {
		return "http"
	}
	return "https"
}
