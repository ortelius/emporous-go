package orasclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatherDescriptors(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		expDigest := "sha256:2e30f6131ce2164ed5ef017845130727291417d60a1be6fad669bdc4473289cd"

		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors("", "testdata/workspace/fish.jpg")
		require.NoError(t, err)
		require.Len(t, desc, 1)
		require.Equal(t, expDigest, desc[0].Digest.String())
	})
}

// TODO(jpower432): Create a mock client to mock non-tested actions
func TestGenerateManifest(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		expDigest := "sha256:792a4e91200098b062d4bf3c8b95dc11c7c20d3a18d8208cd00cae44f6147e37"

		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.GatherDescriptors("", "testdata/workspace/fish.jpg")
		require.NoError(t, err)
		configDesc, err := c.GenerateConfig(nil)
		require.NoError(t, err)
		mdesc, err := c.GenerateManifest("localhost:5000/test:latest", configDesc, nil, desc...)
		require.NoError(t, err)
		require.Equal(t, expDigest, mdesc.Digest.String())
	})
}
