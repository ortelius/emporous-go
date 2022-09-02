package render

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/uor-framework/uor-client-go/cli/log"
	"github.com/uor-framework/uor-client-go/cli/options"
)

func TestRenderComplete(t *testing.T) {
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
			args: []string{"./testdata"},
			expOpts: &Options{
				Output:  "client-workspace",
				RootDir: "./testdata",
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

func TestRenderValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *Options
		expError string
	}

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &Options{
				RootDir: "./testdata",
			},
		},
		{
			name: "Invalid/RootDirDoesNotExist",
			opts: &Options{
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

func TestRenderRun(t *testing.T) {

	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	type spec struct {
		name     string
		opts     *Options
		expected string
		expError string
	}

	cases := []spec{
		{
			name: "Success/FlatWorkspace",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "./testdata/flatworkspace",
			},
			expected: "./testdata/expected/flatworkspace",
		},
		{
			name: "Success/MultiLevelWorkspace",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "./testdata/multi-level-workspace",
			},
			expected: "./testdata/expected/multi-level-workspace",
		},
		{
			name: "Success/UORParsing",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "./testdata/uor-template",
			},
			expected: "./testdata/expected/uor-template",
		},
		{
			name: "Failure/TwoRoots",
			opts: &Options{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				RootDir: "./testdata/tworoots",
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
				actual := walkDir(t, c.opts.Output)
				expected := walkDir(t, c.expected)

				for path, data1 := range actual {
					t.Log("checking path " + path)
					data2, ok := expected[path]
					require.True(t, ok)
					require.Equal(t, digest.FromBytes(data2).String(), digest.FromBytes(data1).String())
				}
			}
		})
	}
}

func walkDir(t *testing.T, dir string) map[string][]byte {
	files := map[string][]byte{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		files[filepath.Base(path)] = data

		return nil
	})
	require.NoError(t, err)
	return files
}
