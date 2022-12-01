package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAggregateComplete(t *testing.T) {
	type spec struct {
		name       string
		args       []string
		opts       *AggregateOptions
		assertFunc func(config *AggregateOptions) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com", "myquery.yaml"},
			assertFunc: func(config *AggregateOptions) bool {
				return config.AttributeQuery == "myquery.yaml" && config.RegistryHost == "test-registry.com"
			},
			opts: &AggregateOptions{
				CreateOptions: &CreateOptions{},
			},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			opts:     &AggregateOptions{},
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
