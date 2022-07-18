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
		ctx := context.TODO()
		expDigest := "sha256:2e30f6131ce2164ed5ef017845130727291417d60a1be6fad669bdc4473289cd"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors(ctx, "", testdata)
		require.NoError(t, err)
		require.Len(t, desc, 1)
		require.Equal(t, expDigest, desc[0].Digest.String())
	})
}

// TODO(jpower432): Create a mock client to mock non-tested actions
func TestGenerateManifest(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		ctx := context.TODO()
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(ctx, []byte("{}"), nil)
		require.NoError(t, err)
		mdesc, err := c.GenerateManifest(ctx, "localhost:5000/test:latest", configDesc, nil, desc...)
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
	notExistTag := "latest"
	notExistRef := fmt.Sprintf("%s/notexist:%s", u.Host, notExistTag)
	images := []string{fmt.Sprintf("%s/test:latest", u.Host), fmt.Sprintf("%s/test2:latest", u.Host)}
	testdata := filepath.Join("testdata", "workspace", "fish.jpg")

	ctx := context.TODO()

	t.Run("Success/PushOneImage", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.GatherDescriptors(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(ctx, []byte("{}"), nil)
		require.NoError(t, err)

		mdesc, err := c.GenerateManifest(ctx, ref, configDesc, nil, descs...)
		require.NoError(t, err)
		desc, err := c.Execute(context.TODO(), ref, registryclient.TypePush)
		require.NoError(t, err)
		require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		require.Equal(t, expDigest, desc.Digest.String())
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullOneImage", func(t *testing.T) {
		tmp := t.TempDir()
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		c, err := NewClient(WithPlainHTTP(true), WithOutputDir(tmp))
		require.NoError(t, err)
		desc, err := c.Execute(context.TODO(), ref, registryclient.TypePull)
		require.NoError(t, err)
		require.Equal(t, expDigest, desc.Digest.String())
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PushMultipleImages", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.GatherDescriptors(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(ctx, []byte("{}"), nil)
		require.NoError(t, err)

		for _, ref := range images {
			mdesc, err := c.GenerateManifest(ctx, ref, configDesc, nil, descs...)
			require.NoError(t, err)
			desc, err := c.Execute(context.TODO(), ref, registryclient.TypePush)
			require.NoError(t, err)
			require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		}
		require.NoError(t, c.Destroy())
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
		require.NoError(t, c.Destroy())
	})

	t.Run("Failure/UnsupportedAction", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), "localhost:5001/fail", 3)
		require.EqualError(t, err, "unsupported action type")
		require.NoError(t, c.Destroy())
	})

	t.Run("Failure/InvalidAction", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), "localhost:5001/fail", registryclient.TypeInvalid)
		require.EqualError(t, err, "action type must be set")
		require.NoError(t, c.Destroy())
	})

	t.Run("Failure/ImageDoesNotExist", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, err = c.Execute(context.TODO(), notExistRef, registryclient.TypePull)
		require.EqualError(t, err, fmt.Sprintf("%s: not found", notExistTag))
		require.NoError(t, c.Destroy())
	})
}
