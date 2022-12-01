package queries

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestResolveQuery(t *testing.T) {
	manifest := []byte("hello world")
	digest := digest.FromBytes(manifest)
	server := httptest.NewServer(testutils.NewRegistry(t, nil, [][]byte{manifest}))
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := auth.DefaultClient

	queryFN := QueryParamsFn(func(values url.Values) {
		values.Add("attributes", `{"test": "me}`)
	})

	results, err := ResolveQuery(context.Background(), u.Host, queryFN, client, true)
	require.NoError(t, err)

	var index ocispec.Index
	err = json.Unmarshal(results, &index)
	require.NoError(t, err)

	require.Equal(t, []ocispec.Descriptor{{MediaType: "application/vnd.oci.artifact.manifest.v1+json", Digest: digest, Size: int64(len(manifest))}}, index.Manifests)
}
