package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/content/layout"
)

func TestPushComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *PushOptions
		expOpts  *PushOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com/image:latest"},
			expOpts: &PushOptions{
				Destination: "test-registry.com/image:latest",
			},
			opts: &PushOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PushOptions{},
			opts:     &PushOptions{},
			expError: "bug: expecting one argument",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Complete(c.args)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expOpts, c.opts)
			}
		})
	}
}

func TestPushRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *PushOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/Stored",
			opts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/success:latest", u.Host),
				PlainHTTP:   true,
			},
		},
		{
			name: "Failure/NotStored",
			opts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger:   testlogr,
					cacheDir: "testdata/cache",
				},
				Destination: "locahost:5001/client-flat-test:latest",
				PlainHTTP:   true,
			},
			expError: "error publishing content to locahost:5001/client-flat-test:latest:" +
				" descriptor for reference locahost:5001/client-flat-test:latest is not stored",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))

			if c.opts.cacheDir == "" {
				c.opts.cacheDir = cache
				prepCache(t, c.opts.Destination, cache)
			}

			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// prepCache will push a hello.txt artifact into the
// cache for retrieval. Uses methods from oras-go.
func prepCache(t *testing.T, ref string, cacheDir string) {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	ctx := context.TODO()

	ociStore, err := layout.New(cacheDir)
	require.NoError(t, err)
	layerDesc, err := pushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, ociStore)
	require.NoError(t, err)
	if layerDesc.Annotations == nil {
		layerDesc.Annotations = map[string]string{}
	}
	layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

	config := []byte("{}")
	configDesc, err := pushBlob(ctx, ocispec.MediaTypeImageConfig, config, ociStore)
	require.NoError(t, err)

	manifest, err := generateManifest(configDesc, layerDesc)
	require.NoError(t, err)

	manifestDesc, err := pushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, ociStore)
	require.NoError(t, err)

	require.NoError(t, ociStore.Tag(ctx, manifestDesc, ref))

}
