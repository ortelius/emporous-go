package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatExamples(t *testing.T) {
	type spec struct {
		name     string
		examples []Example
		expRes   string
	}

	cases := []spec{
		{
			name: "Success/NoExamples",
		},
		{
			name: "Success/OneExample",
			examples: []Example{
				{
					RootCommand:   "test",
					CommandString: "subcommand",
					Descriptions:  []string{"This is a test."},
				},
			},
			expRes: "  # This is a test.\n  test subcommand",
		},
		{
			name: "Success/OneExampleMultipleDescriptionLines",
			examples: []Example{
				{
					RootCommand:   "test",
					CommandString: "subcommand",
					Descriptions: []string{
						"This is a test.",
						"The default is false",
					},
				},
			},
			expRes: "  # This is a test.\n  # The default is false\n  test subcommand",
		},
		{
			name: "Success/TwoExamples",
			examples: []Example{
				{
					RootCommand:   "test",
					CommandString: "subcommand",
					Descriptions:  []string{"This is a test."},
				},
				{
					RootCommand:   "test",
					CommandString: "subcommand --flag",
					Descriptions:  []string{"This is a test with a flag."},
				},
			},
			expRes: "  # This is a test.\n  test subcommand\n  \n  # This is a test with a flag.\n  test subcommand --flag",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.expRes, FormatExamples(c.examples...))
		})
	}

}
