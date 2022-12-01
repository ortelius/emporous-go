package commands

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
)

// CreateOptions describe configuration options that can
// be set using the create subcommand.
type CreateOptions struct {
	*options.Common
	options.Remote
	options.RemoteAuth
}

// NewCreateCmd creates a new cobra.Command for the create subcommand.
func NewCreateCmd(common *options.Common) *cobra.Command {
	o := CreateOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "create",
		Short:         "Create artifacts from existing OCI artifacts",
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	o.Remote.BindFlags(cmd.PersistentFlags())
	o.RemoteAuth.BindFlags(cmd.PersistentFlags())

	cmd.AddCommand(NewAggregateCmd(&o))

	return cmd
}
