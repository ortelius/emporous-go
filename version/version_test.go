package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteVersion(t *testing.T) {
	type spec struct {
		name          string
		testVersion   string
		testCommit    string
		testDate      string
		testBuildData string
		expError      string
		assertFunc    func(string) bool
	}

	cases := []spec{
		{
			name: "Valid/NoVariablesSet",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.0-unknown")
			},
		},
		{
			name:        "Valid/VariablesSet",
			testVersion: "v0.0.1",
			testCommit:  "commit",
			testDate:    "today",
			assertFunc: func(s string) bool {
				return strings.Contains(s, "v0.0.1") && strings.Contains(s, "commit") && strings.Contains(s, "today")
			},
		},
		{
			name:          "Valid/VariablesSetWithBuildData",
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
			if c.testVersion != "" {
				version = c.testVersion
			}
			buildDate = c.testDate
			commit = c.testCommit
			buildData = c.testBuildData
			err := WriteVersion(out)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(out.String()))
			}
		})
	}
}
