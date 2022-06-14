package orasclient

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/registryclient"
)

func TestGatherDescriptors(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		expDigest := "sha256:2e30f6131ce2164ed5ef017845130727291417d60a1be6fad669bdc4473289cd"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors("", testdata)
		require.NoError(t, err)
		require.Len(t, desc, 1)
		require.Equal(t, expDigest, desc[0].Digest.String())
	})
}

// TODO(jpower432): Create a mock client to mock non-tested actions
func TestGenerateManifest(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		expDigest := "sha256:792a4e91200098b062d4bf3c8b95dc11c7c20d3a18d8208cd00cae44f6147e37"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors("", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(nil)
		require.NoError(t, err)
		mdesc, err := c.GenerateManifest("localhost:5000/test:latest", configDesc, nil, desc...)
		require.NoError(t, err)
		require.Equal(t, expDigest, mdesc.Digest.String())
	})
}

func TestExecute(t *testing.T) {
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	ref := fmt.Sprintf("%s/test:latest", u.Host)
	notExistRef := fmt.Sprintf("%s/notexist:latest", u.Host)
	images := []string{fmt.Sprintf("%s/test:latest", u.Host), fmt.Sprintf("%s/test2:latest", u.Host)}
	testdata := filepath.Join("testdata", "workspace", "fish.jpg")

	t.Run("Success/PushOneImage", func(t *testing.T) {
		expDigest := "sha256:792a4e91200098b062d4bf3c8b95dc11c7c20d3a18d8208cd00cae44f6147e37"
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.GatherDescriptors("", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(nil)
		require.NoError(t, err)

		mdesc, err := c.GenerateManifest(ref, configDesc, nil, descs...)
		require.NoError(t, err)
		desc, err := c.Execute(context.TODO(), ref, registryclient.TypePush)
		require.NoError(t, err)
		require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		require.Equal(t, expDigest, desc.Digest.String())

	})

	t.Run("Success/PullOneImage", func(t *testing.T) {
		tmp := t.TempDir()
		expDigest := "sha256:792a4e91200098b062d4bf3c8b95dc11c7c20d3a18d8208cd00cae44f6147e37"
		c, err := NewClient(WithPlainHTTP(true), WithOutputDir(tmp))
		require.NoError(t, err)
		desc, err := c.Execute(context.TODO(), ref, registryclient.TypePull)
		require.NoError(t, err)
		require.Equal(t, expDigest, desc.Digest.String())
	})

	t.Run("Success/PushMultipleImages", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.GatherDescriptors("", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(nil)
		require.NoError(t, err)

		for _, ref := range images {
			mdesc, err := c.GenerateManifest(ref, configDesc, nil, descs...)
			require.NoError(t, err)
			desc, err := c.Execute(context.TODO(), ref, registryclient.TypePush)
			require.NoError(t, err)
			require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		}
	})

	t.Run("Success/PullMultipleImages", func(t *testing.T) {
		tmp := t.TempDir()
		c, err := NewClient(WithPlainHTTP(true), WithOutputDir(tmp))
		require.NoError(t, err)
		for _, ref := range images {
			_, err := c.Execute(context.TODO(), ref, registryclient.TypePull)
			require.NoError(t, err)
			_, err = os.Stat(filepath.Join(tmp, testdata))
			require.NoError(t, err)
		}
	})

	t.Run("Failure/UnsupportedAction", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), "localhost:5001/fail", 3)
		require.EqualError(t, err, "unsupported action type")

	})

	t.Run("Failure/InvalidAction", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), "localhost:5001/fail", registryclient.TypeInvalid)
		require.EqualError(t, err, "action type must be set")
	})

	t.Run("Failure/ImageDoesNotExist", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), notExistRef, registryclient.TypePull)
		require.EqualError(t, err, fmt.Sprintf("%s: not found", notExistRef))
	})
}
