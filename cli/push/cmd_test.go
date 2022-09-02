package push

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

	"github.com/uor-framework/uor-client-go/cli/internal/testutils"
	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/cli/options"
)

func TestPushComplete(t *testing.T) {
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
				Destination: "test-registry.com/image:latest",
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

func TestPushRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *Options
		expError string
	}

	cases := []spec{
		{
			name: "Success/Stored",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/success:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
		},
		{
			name: "Failure/NotStored",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger:   testlogr,
					CacheDir: "./testdata",
				},
				Destination: "localhost:5001/client-flat-test:latest",
				Remote: options.Remote{
					PlainHTTP: true,
				},
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
				err := testutils.PrepCache(c.opts.Destination, cache, nil)
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
