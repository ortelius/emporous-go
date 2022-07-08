package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/cli/log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
			args: []string{"testdata", "test-registry.com/image:latest"},
			expOpts: &PushOptions{
				RootDir:     "testdata",
				Destination: "test-registry.com/image:latest",
			},
			opts: &PushOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PushOptions{},
			opts:     &PushOptions{},
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

func TestPushValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *PushOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &PushOptions{
				RootDir: "testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &PushOptions{
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
			name: "Success/FlatWorkspace",
			opts: &PushOptions{
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
				PlainHTTP:   true,
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &PushOptions{
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
				PlainHTTP:   true,
			},
		},
		{
			name: "SuccessTwoRoots",
			opts: &PushOptions{
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
				PlainHTTP:   true,
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
				// TODO(jpower432): check image is pullable
				// Will do after adding pulling functionality
			}
		})
	}
}