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
	"github.com/uor-framework/uor-client-go/cli/log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestBuildComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *BuildOptions
		expOpts  *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"testdata", "test-registry.com/image:latest"},
			expOpts: &BuildOptions{
				RootDir:     "testdata",
				Destination: "test-registry.com/image:latest",
			},
			opts: &BuildOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &BuildOptions{},
			opts:     &BuildOptions{},
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

func TestBuildValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &BuildOptions{
				RootDir: "testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &BuildOptions{
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

func TestBuildRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
				RootDir:     "testdata/flatworkspace",
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
				RootDir:     "testdata/multi-level-workspace",
			},
		},
		{
			name: "SuccessTwoRoots",
			opts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-tworoots-test:latest", u.Host),
				RootDir:     "testdata/tworoots",
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
			}
		})
	}
}
