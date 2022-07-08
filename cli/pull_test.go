package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/uor-framework/client/cli/log"
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
		name      string
		opts      *PullOptions
		fileExist bool
		expError  string
	}

	cases := []spec{
		{
			name: "Success/NoAttributes",
			opts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source:    fmt.Sprintf("%s/client-test:latest", u.Host),
				PlainHTTP: true,
			},
			fileExist: true,
		},
		{
			name: "Success/OneMatchingAnnotation",
			opts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-test:latest", u.Host),
				Attributes: map[string]string{
					"test": "annotation",
				},
				PlainHTTP: true,
			},
			fileExist: true,
		},
		{
			name: "Success/NoMatchingAnnotation",
			opts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-test:latest", u.Host),
				Attributes: map[string]string{
					"test2": "annotation",
				},
				PlainHTTP: true,
			},
			fileExist: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tmp := t.TempDir()
			c.opts.Output = tmp
			prepTestArtifact(t, c.opts.Source)
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				actual := filepath.Join(tmp, "hello.txt")
				_, err = os.Stat(actual)
				if c.fileExist {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, os.ErrNotExist)
				}
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
	ctx := context.TODO()
	// Push file(s) w custom mediatype to registry
	memoryStore := memory.New()
	layerDesc, err := pushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, memoryStore)
	require.NoError(t, err)
	if layerDesc.Annotations == nil {
		layerDesc.Annotations = map[string]string{}
	}
	layerDesc.Annotations["test"] = "annotation"
	layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

	config := []byte("{}")
	configDesc, err := pushBlob(ctx, ocispec.MediaTypeImageConfig, config, memoryStore)
	require.NoError(t, err)

	manifest, err := generateManifest(configDesc, layerDesc)
	require.NoError(t, err)

	manifestDesc, err := pushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, memoryStore)
	require.NoError(t, err)

	require.NoError(t, memoryStore.Tag(ctx, manifestDesc, ref))

	repo, err := remote.NewRepository(ref)
	require.NoError(t, err)
	repo.PlainHTTP = true
	_, err = oras.Copy(context.TODO(), memoryStore, ref, repo, "", oras.DefaultCopyOptions)
	require.NoError(t, err)
}

func pushBlob(ctx context.Context, mediaType string, blob []byte, target oras.Target) (ocispec.Descriptor, error) {
	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    digest.FromBytes(blob),
		Size:      int64(len(blob)),
	}
	return desc, target.Push(ctx, desc, bytes.NewReader(blob))
}

func generateManifest(configDesc ocispec.Descriptor, layers ...ocispec.Descriptor) ([]byte, error) {
	manifest := ocispec.Manifest{
		Config:    configDesc,
		Layers:    layers,
		Versioned: specs.Versioned{SchemaVersion: 2},
	}
	return json.Marshal(manifest)
}
