package cli

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"
)

type RootOptions struct {
	genericclioptions.IOStreams
}

var clientLong = templates.LongDesc(
	`
	This client helps you build and deploy sets of OCI artifacts to use
	with existing clients.
	`,
)

func NewRootCmd() *cobra.Command {
	o := RootOptions{}
	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           os.Args[0],
		Short:         "UOR Client",
		Long:          clientLong,
		SilenceErrors: false,
		SilenceUsage:  false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewBuildCmd(&o))

	return cmd
}
