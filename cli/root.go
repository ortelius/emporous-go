package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/cli/log"
)

type RootOptions struct {
	IOStreams genericclioptions.IOStreams
	LogLevel  string
	Logger    log.Logger
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
		Use:           filepath.Base(os.Args[0]),
		Short:         "UOR Client",
		Long:          clientLong,
		SilenceErrors: false,
		SilenceUsage:  false,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			logger, err := log.NewLogger(o.IOStreams.Out, o.LogLevel)
			if err != nil {
				return err
			}
			o.Logger = logger
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&o.LogLevel, "loglevel", "l", "info",
		"Log level (debug, info, warn, error, fatal)")

	cmd.AddCommand(NewBuildCmd(&o))
	cmd.AddCommand(NewPushCmd(&o))

	return cmd
}
