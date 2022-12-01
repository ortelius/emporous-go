package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	orasregistry "oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/log"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
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
			args: []string{"test-registry.com/image:latest"},
			expOpts: &PullOptions{
				Source: "test-registry.com/image:latest",
				Output: ".",
			},
			opts: &PullOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PullOptions{},
			opts:     &PullOptions{},
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
	testlogr, err := log.NewLogrusLogger(io.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name       string
		opts       *PullOptions
		assertFunc func(string) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Success/NoAttributes",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:   fmt.Sprintf("%s/client-test:latest", u.Host),
				NoVerify: true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/PullAll",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:   fmt.Sprintf("%s/client-linked:latest", u.Host),
				PullAll:  true,
				NoVerify: true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				if err != nil {
					return false
				}
				actual = filepath.Join(path, "aggregate.txt")
				_, err = os.Stat(actual)
				if err != nil {
					return false
				}
				actual = filepath.Join(path, "aggregate2.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/PullAllWithAttributes",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:         fmt.Sprintf("%s/client-linked-attr:latest", u.Host),
				PullAll:        true,
				AttributeQuery: "testdata/configs/link.yaml",
				NoVerify:       true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "aggregate.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/OneMatchingAnnotation",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "testdata/configs/match.yaml",
				NoVerify:       true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/NoMatchingAnnotation",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "testdata/configs/nomatch.yaml",
				NoVerify:       true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return errors.Is(err, os.ErrNotExist)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tmp := t.TempDir()
			c.opts.Output = tmp
			prepTestArtifact(t, c.opts.Source)

			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.opts.CacheDir = cache

			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(tmp))
			}
		})
	}
}

// prepTestArtifact will push a hello.txt artifact into the
// registry for retrieval. Uses methods from oras-go.
func prepTestArtifact(t *testing.T, ref string) {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")

	aggregateRef := fmt.Sprintf("%s-aggregate", ref)
	aggregateDesc := prepLinks(t, aggregateRef)

	aggregationJSON, err := json.Marshal(aggregateDesc)
	require.NoError(t, err)

	manifestAnnotations := map[string]string{
		uorspec.AnnotationLink: string(aggregationJSON),
	}
	_, err = publishFunc(fileName, ref, fileContent, map[string]string{"test": "annotation"}, manifestAnnotations)
	require.NoError(t, err)
}

// prepLinks will push links into the
// registry for retrieval. Uses methods from oras-go.
func prepLinks(t *testing.T, ref string) []ocispec.Descriptor {
	fileName := "aggregate.txt"
	fileContent := []byte("Hello Again World!\n")
	ref1 := fmt.Sprintf("%s-ref1", ref)

	r, err := orasregistry.ParseReference(ref)
	require.NoError(t, err)
	linkAttr := descriptor.Properties{
		Link: &uorspec.LinkAttributes{
			RegistryHint:  r.Registry,
			NamespaceHint: r.Repository,
			Transitive:    true,
		},
	}
	linkJSON, err := json.Marshal(linkAttr)
	require.NoError(t, err)
	ref1Annotations := map[string]string{
		uorspec.AnnotationUORAttributes: string(linkJSON),
	}
	desc1, err := publishFunc(fileName, ref1, fileContent, map[string]string{"test": "linkedannotation"}, ref1Annotations)
	require.NoError(t, err)
	fileName2 := "aggregate2.txt"
	fileContent2 := []byte("Hello Again Again World !\n")
	ref2 := fmt.Sprintf("%s-ref2", ref)
	ref2Annotations := map[string]string{
		uorspec.AnnotationUORAttributes: string(linkJSON),
	}
	desc2, err := publishFunc(fileName2, ref2, fileContent2, map[string]string{"test": "annotation"}, ref2Annotations)
	require.NoError(t, err)

	return []ocispec.Descriptor{desc1, desc2}
}

func publishFunc(fileName, ref string, fileContent []byte, layerAnnotations, manifestAnnotations map[string]string) (ocispec.Descriptor, error) {
	ctx := context.TODO()
	// Push file(s) w custom mediatype to registry
	memoryStore := memory.New()
	layerDesc, err := pushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	layerDesc.Annotations = layerAnnotations
	if layerDesc.Annotations == nil {
		layerDesc.Annotations = map[string]string{}
	}
	layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

	config := []byte("{}")
	configDesc, err := pushBlob(ctx, ocispec.MediaTypeImageConfig, config, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	manifest, err := generateManifest(configDesc, manifestAnnotations, layerDesc)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	manifestDesc, err := pushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	err = memoryStore.Tag(ctx, manifestDesc, ref)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	repo, err := remote.NewRepository(ref)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}
	repo.PlainHTTP = true
	return oras.Copy(context.TODO(), memoryStore, ref, repo, "", oras.DefaultCopyOptions)
}
