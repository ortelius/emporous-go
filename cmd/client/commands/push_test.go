package commands

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

	"github.com/emporous/emporous-go/cmd/client/commands/options"
	"github.com/emporous/emporous-go/log"
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
			args: []string{"test-registry.com/image:latest"},
			expOpts: &PushOptions{
				Destination: "test-registry.com/image:latest",
			},
			opts: &PushOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PushOptions{},
			opts:     &PushOptions{},
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

func TestPushRun(t *testing.T) {
	testlogr, err := log.NewLogrusLogger(ioutil.Discard, "debug")
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
			name: "Success/Stored",
			opts: &PushOptions{
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
				Destination: fmt.Sprintf("%s/success:latest", u.Host),
			},
		},
		{
			name: "Failure/NotStored",
			opts: &PushOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger:   testlogr,
					CacheDir: "testdata/cache",
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Destination: "localhost:5001/client-flat-test:latest",
			},
			expError: "error publishing content to localhost:5001/client-flat-test:latest:" +
				" descriptor for reference localhost:5001/client-flat-test:latest is not stored",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))

			if c.opts.CacheDir == "" {
				c.opts.CacheDir = cache
				err := prepCache(c.opts.Destination, cache, nil)
				require.NoError(t, err)
			}

			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
