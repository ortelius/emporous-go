package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateInventoryComplete(t *testing.T) {
	type spec struct {
		name       string
		args       []string
		opts       *InventoryOptions
		assertFunc func(config *InventoryOptions) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com/test:latest"},
			assertFunc: func(config *InventoryOptions) bool {
				return config.Format == "spdx22json" && config.Source == "test-registry.com/test:latest"
			},
			opts: &InventoryOptions{
				CreateOptions: &CreateOptions{},
			},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			opts:     &InventoryOptions{},
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
				require.True(t, c.assertFunc(c.opts))
			}
		})
	}
}
