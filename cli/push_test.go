package cli

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/cli/log"
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
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *PushOptions
		expError string
	}

	cases := []spec{
		{
			name: "Failure/NotStored",
			opts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger:   testlogr,
					cacheDir: "testdata/cache",
				},
				Destination: "locahost:5001/client-flat-test:latest",
				PlainHTTP:   true,
			},
			expError: "error publishing content to locahost:5001/client-flat-test:latest:" +
				" descriptor for reference locahost:5001/client-flat-test:latest is not stored",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
