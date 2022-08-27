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
	"text/template"

	"github.com/google/go-containerregistry/pkg/registry"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/ocimanifest"
)

func TestBuildCollectionComplete(t *testing.T) {
	type spec struct {
		name       string
		args       []string
		opts       *BuildCollectionOptions
		assertFunc func(config *BuildCollectionOptions) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"testdata", "test-registry.com/image:latest"},
			assertFunc: func(config *BuildCollectionOptions) bool {
				return config.RootDir == "testdata" && config.Destination == "test-registry.com/image:latest"
			},
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{},
			},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			opts:     &BuildCollectionOptions{},
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
				require.True(t, c.assertFunc(c.opts))
			}
		})
	}
}

func TestBuildCollectionValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *BuildCollectionOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &BuildCollectionOptions{
				RootDir: "testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &BuildCollectionOptions{
				RootDir: "fake",
			},
			expError: "workspace directory \"fake\": stat fake: no such file or directory",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Validate()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBuildCollectionRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *BuildCollectionOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				RootDir: "testdata/flatworkspace",
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				RootDir: "testdata/multi-level-workspace",
			},
		},
		{
			name: "SuccessTwoRoots",

			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-tworoot-test:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				RootDir: "testdata/tworoots",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.opts.cacheDir = cache
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				_, err := os.Stat(filepath.Join(c.opts.cacheDir, "index.json"))
				require.NoError(t, err)
			}
		})
	}
}

func TestBuildCollectionRun_WithDSConfig(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *BuildCollectionOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/BasicConfig",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-test1:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				DSConfig:  "testdata/configs/dataset-config-basic.yaml",
				RootDir:   "testdata/multi-level-workspace",
				PlainHTTP: true,
			},
		},
		{
			name: "Success/WithSchema",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-test2:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				DSConfig:  "testdata/configs/dataset-config-schema.yaml",
				RootDir:   "testdata/multi-level-workspace",
				PlainHTTP: true,
			},
		},
		{
			name: "Success/WithLinks",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-test3:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				DSConfig:  "testdata/configs/dataset-config-links.yaml",
				RootDir:   "testdata/multi-level-workspace",
				PlainHTTP: true,
			},
		},
		{
			name: "Failure/LinksHasNoSchema",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-test3:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				DSConfig:  "testdata/configs/dataset-config-invalidlinks.yaml",
				RootDir:   "testdata/multi-level-workspace",
				PlainHTTP: true,
			},
			expError: fmt.Sprintf("collection \"%s/schema-test:latest\": no schema", u.Host),
		},
		{
			name: "Failure/InvalidSchema",
			opts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-test3:latest", u.Host),
					RootOptions: &RootOptions{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				DSConfig:  "testdata/configs/dataset-config-invalidschema.yaml",
				RootDir:   "testdata/multi-level-workspace",
				PlainHTTP: true,
			},
			expError: fmt.Sprintf("reference %s/test:latest is not a schema address", u.Host),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.opts.cacheDir = cache

			templateValues := prepCollectionArtifacts(t, u.Host)
			initialConfig, err := ioutil.ReadFile(c.opts.DSConfig)
			require.NoError(t, err)
			tpl, err := template.New(c.name).Parse(string(initialConfig))
			require.NoError(t, err)
			finalConfigPath := filepath.Join(t.TempDir(), "test.yaml")
			finalConfig, err := os.Create(finalConfigPath)
			require.NoError(t, err)
			require.NoError(t, tpl.Execute(finalConfig, templateValues))
			require.NoError(t, finalConfig.Close())

			c.opts.DSConfig = finalConfigPath

			err = c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				_, err := os.Stat(filepath.Join(c.opts.cacheDir, "index.json"))
				require.NoError(t, err)
			}
		})
	}
}

// prepCollectionsArtifact pushes a test schema and test collection for testing.
// Uses methods from oras-go. It returns the references and the corresponding template values.
func prepCollectionArtifacts(t *testing.T, host string) map[string]string {
	publishFunc := func(fileName, ref, layerMediaType string, fileContent []byte, layerAnnotations, manifestAnnotations map[string]string) {
		ctx := context.TODO()
		// Push file(s) w custom mediatype to registry
		memoryStore := memory.New()
		layerDesc, err := pushBlob(ctx, layerMediaType, fileContent, memoryStore)
		require.NoError(t, err)
		layerDesc.Annotations = layerAnnotations
		if layerDesc.Annotations == nil {
			layerDesc.Annotations = map[string]string{}
		}
		layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

		config := []byte("{}")
		configDesc, err := pushBlob(ctx, ocispec.MediaTypeImageConfig, config, memoryStore)
		require.NoError(t, err)

		manifest, err := generateManifest(configDesc, manifestAnnotations, layerDesc)
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

	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	testCollection := fmt.Sprintf("%s/test:latest", host)
	testCollectionAnnotations := map[string]string{
		ocimanifest.AnnotationSchema: "test.com/schema:latest",
	}
	publishFunc(fileName, testCollection, ocispec.MediaTypeImageLayer, fileContent, map[string]string{"test": "annotation"}, testCollectionAnnotations)

	schemaName := "schema"
	schemaContent := []byte("{\"type\":\"object\",\"properties\":{\"test\":{\"type\": \"string\"}},\"required\":[\"test\"]}")
	schemaRef := fmt.Sprintf("%s/schema-test:latest", host)
	publishFunc(schemaName, schemaRef, ocimanifest.UORSchemaMediaType, schemaContent, nil, nil)

	return map[string]string{
		"linkedCollection": testCollection,
		"schemaAddress":    schemaRef,
	}
}
