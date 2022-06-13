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
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/cli/log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

func TestPullComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *PullOptions
		expOpts  *PullOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com/image:latest", "test"},
			expOpts: &PullOptions{
				Source: "test-registry.com/image:latest",
				Output: "test",
			},
			opts: &PullOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PullOptions{},
			opts:     &PullOptions{},
			expError: "bug: expecting two arguments",
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

func TestPullValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *PullOptions
		expError string
	}

	tmp := t.TempDir()

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &PullOptions{
				Output: "testdata",
			},
		},
		{
			name: "Valid/RootDirDoesNotExist",
			opts: &PullOptions{
				Output: filepath.Join(tmp, "fake"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Validate()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				_, err = os.Stat(c.opts.Output)
				require.NoError(t, err)
			}
		})
	}
}

func TestPullRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *PullOptions
		expError string
	}

	tmp := t.TempDir()

	cases := []spec{
		{
			name: "SuccessOneImage",
			opts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-tworoots-test:latest", u.Host),
				Output: tmp,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			prepTestArtifact(t, c.opts.Source)
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				actual := filepath.Join(tmp, "hello.txt")
				_, err = os.Stat(actual)
				require.NoError(t, err)
			}
		})
	}
}

// prepTestArtifact will push a hello.txt artifact into the
// registry for retrieval. Uses methods from oras-go.
// FIXME(jpower432): Possibly set this up to mirror from files.
func prepTestArtifact(t *testing.T, ref string) {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	// Push file(s) w custom mediatype to registry
	memoryStore := content.NewMemory()
	desc, err := memoryStore.Add(fileName, "", fileContent)
	require.NoError(t, err)

	manifest, manifestDesc, config, configDesc, err := content.GenerateManifestAndConfig(nil, nil, desc)
	require.NoError(t, err)
	memoryStore.Set(configDesc, config)
	err = memoryStore.StoreManifest(ref, manifestDesc, manifest)
	require.NoError(t, err)
	registry, err := content.NewRegistry(content.RegistryOptions{PlainHTTP: true})
	require.NoError(t, err)
	desc, err = oras.Copy(context.TODO(), memoryStore, ref, registry, "")
	require.NoError(t, err)
}
