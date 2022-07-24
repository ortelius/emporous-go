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

func TestCLIE2E(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name          string
		pushOpts      *PushOptions
		buildOpts     *BuildOptions
		pullOpts      *PullOptions
		expBuildError string
		expPushError  string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			buildOpts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "testdata/flatworkspace",
			},
			pushOpts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
				PlainHTTP:   true,
			},
			pullOpts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				PlainHTTP: true,
				Output:    t.TempDir(),
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			buildOpts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "testdata/multi-level-workspace",
			},
			pushOpts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
				PlainHTTP:   true,
			},
			pullOpts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Output:    t.TempDir(),
				PlainHTTP: true,
			},
		},
		{
			name: "Success/UORParsing",
			buildOpts: &BuildOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "testdata/uor-template",
			},
			pushOpts: &PushOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-uor-test:latest", u.Host),
				PlainHTTP:   true,
			},
			pullOpts: &PullOptions{
				RootOptions: &RootOptions{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Output:    t.TempDir(),
				PlainHTTP: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.buildOpts.Output = t.TempDir()
			err := c.buildOpts.Run(context.TODO())
			if c.expBuildError != "" {
				require.EqualError(t, err, c.expBuildError)
			} else {
				require.NoError(t, err)
			}

			c.pushOpts.RootDir = c.buildOpts.Output
			err = c.pushOpts.Run(context.TODO())
			if c.expPushError != "" {
				require.EqualError(t, err, c.expBuildError)
			} else {
				require.NoError(t, err)
			}

			c.pullOpts.Source = c.pushOpts.Destination
			err = c.pullOpts.Run(context.TODO())
			if c.expPushError != "" {
				require.EqualError(t, err, c.expBuildError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
