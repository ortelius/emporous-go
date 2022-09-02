package cli

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/uor-client-go/cli/build"
	"github.com/uor-framework/uor-client-go/cli/inspect"
	"github.com/uor-framework/uor-client-go/cli/options"
	"github.com/uor-framework/uor-client-go/cli/pull"
	"github.com/uor-framework/uor-client-go/cli/push"
	"github.com/uor-framework/uor-client-go/cli/render"
	"github.com/uor-framework/uor-client-go/cli/run"
	"github.com/uor-framework/uor-client-go/cli/version"
)

var clientLong = templates.LongDesc(
	`
	The UOR client helps you build, publish, and retrieve UOR collections as an OCI artifact.
	`,
)

// NewClientCmd creates a new cobra.Command for the command root.
func NewClientCmd() *cobra.Command {
	o := options.Common{}
	o.EnvConfig = readEnvConfig()

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
			if err := o.Init(); err != nil {
				return err
			}
			return os.MkdirAll(o.CacheDir, 0750)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	o.BindFlags(cmd.PersistentFlags())

	cmd.AddCommand(render.NewCmd(&o))
	cmd.AddCommand(inspect.NewCmd(&o))
	cmd.AddCommand(build.NewCmd(&o))
	cmd.AddCommand(push.NewCmd(&o))
	cmd.AddCommand(pull.NewCmd(&o))
	cmd.AddCommand(run.NewCmd(&o))
	cmd.AddCommand(version.NewCmd(&o))

	return cmd
}

func readEnvConfig() options.EnvConfig {
	envConfig := options.EnvConfig{}

	devModeString := os.Getenv("UOR_DEV_MODE")
	devMode, err := strconv.ParseBool(devModeString)
	envConfig.UOR_DEV_MODE = err == nil && devMode

	return envConfig
}
