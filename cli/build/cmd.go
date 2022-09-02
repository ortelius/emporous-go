package build

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cli/options"
)

// NewCmd creates a new cobra.Command for the build subcommand.
func NewCmd(commonOpts *options.Common) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "build",
		Short:         "Build and save an OCI artifact from files",
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewSchemaCmd(commonOpts))
	cmd.AddCommand(NewCollectionCmd(commonOpts))

	return cmd
}
