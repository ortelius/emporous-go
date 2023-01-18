package commands

import (
	"github.com/spf13/cobra"

	"github.com/emporous/emporous-go/cmd/client/commands/options"
)

// BuildOptions describe configuration options that can
// be set using the build subcommand.
type BuildOptions struct {
	*options.Common
	Destination string
}

// NewBuildCmd creates a new cobra.Command for the build subcommand.
func NewBuildCmd(common *options.Common) *cobra.Command {
	o := BuildOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "build",
		Short:         "Build and save an OCI artifact from files",
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewBuildSchemaCmd(&o))
	cmd.AddCommand(NewBuildCollectionCmd(&o))

	return cmd
}
