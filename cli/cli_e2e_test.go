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

	"github.com/uor-framework/uor-client-go/cli/build"
	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/cli/options"
	"github.com/uor-framework/uor-client-go/cli/pull"
	"github.com/uor-framework/uor-client-go/cli/push"
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
		pushOpts      *push.Options
		buildOpts     *build.CollectionOptions
		pullOpts      *pull.Options
		expBuildError string
		expPushError  string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			buildOpts: &build.CollectionOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
				RootDir:     "./testdata/flatworkspace",
			},
			pushOpts: &push.Options{
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
			pullOpts: &pull.Options{
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
				Output: t.TempDir(),
			},
		},
		{
			name: "Success/MultiLevelWorkspace",
			buildOpts: &build.CollectionOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote:      options.Remote{},
				Destination: fmt.Sprintf("%s/client-multi-test:latest", u.Host),
				RootDir:     "./testdata/multi-level-workspace",
			},
			pushOpts: &push.Options{
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
			pullOpts: &pull.Options{
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
				Output: t.TempDir(),
			},
		},
		{
			name: "Success/UORWorkspace",
			buildOpts: &build.CollectionOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Destination: fmt.Sprintf("%s/client-uor-test:latest", u.Host),
				RootDir:     "./testdata/uor-template",
			},
			pushOpts: &push.Options{
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
			pullOpts: &pull.Options{
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
				Output: t.TempDir(),
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
