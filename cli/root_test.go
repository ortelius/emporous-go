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

func TestRootComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *RootOptions
		expOpts  *RootOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"testdata", "test-registry.com/client-test:latest"},
			expOpts: &RootOptions{
				Reference: "test-registry.com/client-test:latest",
				RootDir:   "testdata",
			},
			opts: &RootOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{"testdata"},
			expOpts:  &RootOptions{},
			opts:     &RootOptions{},
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

func TestRootValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *RootOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &RootOptions{
				Reference: "test-registry.com/client-test:latest",
				RootDir:   "testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &RootOptions{
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

func TestRootRun(t *testing.T) {
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *RootOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &RootOptions{
				IOStreams: genericclioptions.IOStreams{
					Out:    os.Stdout,
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
				Reference: fmt.Sprintf("%s/client-test:latest", u.Host),
				RootDir:   "testdata/flatworkspace",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
