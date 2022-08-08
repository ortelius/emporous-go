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
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/uor-framework/uor-client-go/cli/log"
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
