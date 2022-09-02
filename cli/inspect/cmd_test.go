package inspect

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/uor-framework/uor-client-go/cli/internal/testutils"
	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/cli/options"
)

func TestInspectValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *Options
		expError string
	}

	cases := []spec{
		{
			name: "Valid/NoInputs",
			opts: &Options{
				Source: "localhost:5001/test:latest",
			},
		},
		{
			name: "Valid/ReferenceOnly",
			opts: &Options{
				Source: "localhost:5001/test:latest",
			},
		},
		{
			name: "Valid/ReferenceAndAttributes",
			opts: &Options{
				Source: "localhost:5001/test:latest",
			},
		},
		{
			name: "Invalid/AttributesOnly",
			opts: &Options{
				AttributeQuery: "notempty",
			},
			expError: "must specify a reference with --reference",
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

func TestInspectRun(t *testing.T) {
	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name        string
		opts        *Options
		annotations map[string]string
		expRes      string
		expError    string
	}

	cases := []spec{
		{
			name: "Success/AttributesMatch",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source:         fmt.Sprintf("%s/success:latest", u.Host),
				AttributeQuery: "./testdata/configs/match.yaml",
			},
			annotations: map[string]string{
				"test": "annotation",
			},
			expRes: "Listing matching descriptors for source:  " + u.Host + "/success:latest\nName" +
				"                                      Digest" +
				"                                                                   Size  MediaType\nhello.txt" +
				"                                 sha256:03ba204e50d126e4674c005e04d82e84c21366780af1f43bd54a37816b6ab340" +
				"  13    application/vnd.oci.image.layer.v1.tar\n",
		},
		{
			name: "Success/NoAttributesMatch",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Source:         fmt.Sprintf("%s/success:latest", u.Host),
				AttributeQuery: "./testdata/configs/nomatch.yaml",
			},
			expRes: "Listing matching descriptors for source:  " + u.Host + "/success:latest\nName" +
				"                                      Digest  Size  MediaType\n",
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
					CacheDir: "./testdata/cache",
				},
				Source: "localhost:5001/client-fake:latest",
			},
			expError: "descriptor for reference localhost:5001/client-fake:latest is not stored",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))

			if c.opts.CacheDir == "" {
				c.opts.CacheDir = cache
				err := testutils.PrepCache(c.opts.Source, cache, c.annotations)
				require.NoError(t, err)
			}

			out := new(strings.Builder)
			c.opts.IOStreams.Out = out

			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expRes, out.String())
			}
		})
	}
}
