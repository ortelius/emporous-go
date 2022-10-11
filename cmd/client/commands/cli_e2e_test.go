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

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/log"
)

func TestCLIE2E(t *testing.T) {
	testlogr, err := log.NewLogrusLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name          string
		pushOpts      *PushOptions
		buildOpts     *BuildCollectionOptions
		pullOpts      *PullOptions
		expBuildError string
		expPushError  string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			buildOpts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Common: &options.Common{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
					Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
				},
				RootDir: "testdata/flatworkspace",
			},
			pushOpts: &PushOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
			pullOpts: &PullOptions{
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
				Output:   t.TempDir(),
				NoVerify: true,
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			buildOpts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Common: &options.Common{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
					Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
				},
				RootDir: "testdata/multi-level-workspace",
			},
			pushOpts: &PushOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
			pullOpts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Output: t.TempDir(),
				Remote: options.Remote{
					PlainHTTP: true,
				},
				NoVerify: true,
			},
		},
		{
			name: "Success/UORParsing",
			buildOpts: &BuildCollectionOptions{
				BuildOptions: &BuildOptions{
					Common: &options.Common{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
					Destination: fmt.Sprintf("%s/client-uor-test:latest", u.Host),
				},
				RootDir: "testdata/uor-template",
			},
			pushOpts: &PushOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-uor-test:latest", u.Host),
				Remote: options.Remote{
					PlainHTTP: true,
				},
			},
			pullOpts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Output: t.TempDir(),
				Remote: options.Remote{
					PlainHTTP: true,
				},
				NoVerify: true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.pushOpts.CacheDir = cache
			c.pullOpts.CacheDir = cache
			c.buildOpts.CacheDir = cache

			err := c.buildOpts.Run(context.TODO())
			if c.expBuildError != "" {
				require.EqualError(t, err, c.expBuildError)
			} else {
				require.NoError(t, err)
			}

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
