package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestGetVersion(t *testing.T) {
	type spec struct {
		name          string
		testVersion   string
		testCommit    string
		testDate      string
		testBuildData string
		opts          *RootOptions
		expError      string
		assertFunc    func(string) bool
	}

	cases := []spec{
		{
			name: "Valid/NoVariablesSet",
			opts: &RootOptions{
				IOStreams: genericclioptions.IOStreams{
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
			},
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.0-unknown")
			},
		},
		{
			name: "Valid/VariablesSet",
			opts: &RootOptions{
				IOStreams: genericclioptions.IOStreams{
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
			},
			testVersion: "v0.0.1",
			testCommit:  "commit",
			testDate:    "today",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.1") && strings.Contains(s, "commit") && strings.Contains(s, "today")
			},
		},
		{
			name: "Valid/VariablesSetWithBuildData",
			opts: &RootOptions{
				IOStreams: genericclioptions.IOStreams{
					In:     os.Stdin,
					ErrOut: os.Stderr,
				},
			},
			testVersion:   "v0.0.1",
			testCommit:    "commit",
			testDate:      "today",
			testBuildData: "dev",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.1+dev")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := new(strings.Builder)
			c.opts.IOStreams.Out = out
			if c.testVersion != "" {
				version = c.testVersion
			}
			buildDate = c.testDate
			commit = c.testCommit
			buildData = c.testBuildData
			err := getVersion(c.opts)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				t.Log(out.String())
				require.True(t, c.assertFunc(out.String()))
			}
		})
	}
}
