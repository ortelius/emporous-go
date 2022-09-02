package pull

import (
	"context"
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
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/uor-framework/uor-client-go/cli/internal/testutils"
	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/cli/options"
	"github.com/uor-framework/uor-client-go/ocimanifest"
)

func TestPullComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *Options
		expOpts  *Options
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com/image:latest"},
			expOpts: &Options{
				Source: "test-registry.com/image:latest",
				Output: ".",
			},
			opts: &Options{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &Options{},
			opts:     &Options{},
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
		opts     *Options
		expError string
	}

	tmp := t.TempDir()

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &Options{
				Output: "../testdata",
			},
		},
		{
			name: "Valid/RootDirDoesNotExist",
			opts: &Options{
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
	testlogr, err := log.NewLogger(io.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name       string
		opts       *Options
		assertFunc func(string) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Success/NoAttributes",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-test:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/PullAll",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-linked:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
				PullAll: true,
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				if err != nil {
					return false
				}
				actual = filepath.Join(path, "hello.linked.txt")
				_, err = os.Stat(actual)
				if err != nil {
					return false
				}
				actual = filepath.Join(path, "hello.linked1.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/PullAllWithAttributes",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source: fmt.Sprintf("%s/client-linked-attr:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
				PullAll:        true,
				AttributeQuery: "./testdata/link.yaml",
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.linked.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/OneMatchingAnnotation",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "./testdata/match.yaml",
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/NoMatchingAnnotation",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "./testdata/nomatch.yaml",
				Remote: options.Remote{
					PlainHTTP: true,
				},
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
			prepTestArtifact(t, c.opts.Source, u.Host)

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
func prepTestArtifact(t *testing.T, ref string, host string) {
	fileName := "hello.txt"
	fileLinkedName := "hello.linked.txt"
	fileLinked1Name := "hello.linked1.txt"
	fileContent := []byte("Hello World!\n")

	publishFunc := func(fileName, ref string, fileContent []byte, layerAnnotations, manifestAnnotations map[string]string) {
		ctx := context.TODO()
		// Push file(s) w custom mediatype to registry
		memoryStore := memory.New()
		layerDesc, err := testutils.PushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, memoryStore)
		require.NoError(t, err)
		if layerDesc.Annotations == nil {
			layerDesc.Annotations = map[string]string{}
		}
		layerDesc.Annotations = layerAnnotations
		layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

		config := []byte("{}")
		configDesc, err := testutils.PushBlob(ctx, ocispec.MediaTypeImageConfig, config, memoryStore)
		require.NoError(t, err)

		manifest, err := testutils.GenerateManifest(configDesc, manifestAnnotations, layerDesc)
		require.NoError(t, err)

		manifestDesc, err := testutils.PushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, memoryStore)
		require.NoError(t, err)

		require.NoError(t, memoryStore.Tag(ctx, manifestDesc, ref))

		repo, err := remote.NewRepository(ref)
		require.NoError(t, err)
		repo.PlainHTTP = true
		_, err = oras.Copy(context.TODO(), memoryStore, ref, repo, "", oras.DefaultCopyOptions)
		require.NoError(t, err)
	}

	linkAnnotations := map[string]string{
		ocimanifest.AnnotationSchema: "test.com/schema:latest",
	}
	linked1Ref := fmt.Sprintf("%s/linked1:test", host)
	publishFunc(fileLinked1Name, linked1Ref, fileContent, map[string]string{"test": "linked1annotation"}, linkAnnotations)
	middleAnnotations := map[string]string{
		ocimanifest.AnnotationSchema:          "test.com/schema:latest",
		ocimanifest.AnnotationSchemaLinks:     "test.com/schema:latest",
		ocimanifest.AnnotationCollectionLinks: linked1Ref,
	}
	linkedRef := fmt.Sprintf("%s/linked:test", host)
	publishFunc(fileLinkedName, linkedRef, fileContent, map[string]string{"test": "linkedannotation"}, middleAnnotations)
	rootAnnotations := map[string]string{
		ocimanifest.AnnotationSchema:          "test.com/schema:latest",
		ocimanifest.AnnotationSchemaLinks:     "test.com/schema:latest",
		ocimanifest.AnnotationCollectionLinks: linkedRef,
	}
	publishFunc(fileName, ref, fileContent, map[string]string{"test": "annotation"}, rootAnnotations)
}
