package cli

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/require"
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
			args: []string{"testdata"},
			expOpts: &BuildOptions{
				Output:  "client-workspace",
				RootDir: "testdata",
			},
			opts: &BuildOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &BuildOptions{},
			opts:     &BuildOptions{},
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
			name: "Valid/DestinationWithPush",
			opts: &BuildOptions{
				Destination: "test-registry.com/client-test:latest",
				RootDir:     "testdata",
				Push:        true,
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &BuildOptions{
				RootDir: "fake",
			},
			expError: "workspace directory \"fake\": stat fake: no such file or directory",
		},
		{
			name: "Invalid/NoReferenceWithPush",
			opts: &BuildOptions{
				RootDir: "testdata",
				Push:    true,
			},
			expError: "destination must be set when using --push",
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
				RootOptions: &RootOptions{IOStreams: genericclioptions.IOStreams{
					Out:    os.Stdout,
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
				},
				Destination: fmt.Sprintf("%s/client-test:latest", u.Host),
				RootDir:     "testdata/flatworkspace",
				Push:        true,
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &BuildOptions{
				RootOptions: &RootOptions{IOStreams: genericclioptions.IOStreams{
					Out:    os.Stdout,
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
				},
				Destination: fmt.Sprintf("%s/client-test:latest", u.Host),
				RootDir:     "testdata/multi-level-workspace",
				Push:        true,
			},
		},
		{
			name: "Success/UORParsing",
			opts: &BuildOptions{
				RootOptions: &RootOptions{IOStreams: genericclioptions.IOStreams{
					Out:    os.Stdout,
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
				},
				Destination: fmt.Sprintf("%s/client-test:latest", u.Host),
				RootDir:     "testdata/uor-template",
				Push:        true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.opts.Output = t.TempDir()
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				// TODO(jpower432): verify resulting this image is now pullable
			}
		})
	}
}
