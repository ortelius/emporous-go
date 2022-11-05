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

func TestBuildSchemaComplete(t *testing.T) {
	type spec struct {
		name       string
		args       []string
		opts       *BuildSchemaOptions
		assertFunc func(config *BuildSchemaOptions) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"testdata/config.yaml", "test-registry.com/image:latest"},
			assertFunc: func(config *BuildSchemaOptions) bool {
				return config.SchemaConfig == "testdata/config.yaml" && config.Destination == "test-registry.com/image:latest"
			},
			opts: &BuildSchemaOptions{
				BuildOptions: &BuildOptions{},
			},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			opts:     &BuildSchemaOptions{},
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
				require.True(t, c.assertFunc(c.opts))
			}
		})
	}
}

func TestBuildSchemaValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *BuildSchemaOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/SchemaExsits",
			opts: &BuildSchemaOptions{
				SchemaConfig: "testdata/configs/schema-config.yaml",
			},
		},
		{
			name: "Invalid/SchemaDoesNotExist",
			opts: &BuildSchemaOptions{
				SchemaConfig: "fake",
			},
			expError: "schema configuration \"fake\": stat fake: no such file or directory",
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

func TestBuildSchemaRun(t *testing.T) {
	testlogr, err := log.NewLogrusLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *BuildSchemaOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &BuildSchemaOptions{
				BuildOptions: &BuildOptions{
					Destination: fmt.Sprintf("%s/client-flat-test:latest", u.Host),
					Common: &options.Common{
						IOStreams: genericclioptions.IOStreams{
							Out:    os.Stdout,
							In:     os.Stdin,
							ErrOut: os.Stderr,
						},
						Logger: testlogr,
					},
				},
				SchemaConfig: "testdata/configs/schema-config.yaml",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.opts.CacheDir = cache
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				_, err := os.Stat(filepath.Join(c.opts.CacheDir, "index.json"))
				require.NoError(t, err)
			}
		})
	}
}
