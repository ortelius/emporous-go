package commands

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/version"
)

// NewVersionCmd creates a new cobra.Command for the version subcommand.
func NewVersionCmd(common *options.Common) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return version.WriteVersion(common.IOStreams.Out)
		},
	}
}
