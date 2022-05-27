package cli

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/cli/log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestBuildComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *BuildOptions
		expOpts  *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"testdata"},
			expOpts: &BuildOptions{
				Output:  "client-workspace",
				RootDir: "testdata",
			},
			opts: &BuildOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &BuildOptions{},
			opts:     &BuildOptions{},
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

func TestBuildValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &BuildOptions{
				RootDir: "testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &BuildOptions{
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

func TestBuildRun(t *testing.T) {

	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *BuildOptions
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &BuildOptions{
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
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &BuildOptions{
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
		},
		{
			name: "Success/UORParsing",
			opts: &BuildOptions{
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
		},
		{
			name: "Failure/TwoRoots",
			opts: &BuildOptions{
				RootOptions: &RootOptions{IOStreams: genericclioptions.IOStreams{
					Out:    os.Stdout,
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
				},
				RootDir: "testdata/tworoots",
			},
			expError: "error building content: error calculating root node: multiple roots found in graph: fish.jpg, fish2.jpg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.opts.Output = t.TempDir()
			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				// Check build artifacts
			}
		})
	}
}
