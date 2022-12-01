package commands

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/emporous/emporous-go/cmd/client/commands/options"
)

var clientLong = templates.LongDesc(
	`
	The Emporous client helps you build, publish, and retrieve Emporous collections as an OCI artifact.

	The workflow to publish a collection is to gather files for a collection in a directory workspace 
	and use the build sub-command. During the build process, the tag for the
	remote destination is specified. 
	
	This build action will store the collection in a build cache. This location can be specified with the EMPOROUS_CACHE environment
	variable. The default location is ~/.emporous/cache. 
	
	After the collection has been stored, it can be retrieved and pushed to the registry with the push sub-command.

	Collections can be retrieved from the cache or the remote location (if not stored) with the pull sub-command. The pull sub-command also
	allows for filtering of the collection with an attribute query configuration file.
	`,
)

// NewRootCmd creates a new cobra.Command for the command root.
func NewRootCmd() *cobra.Command {
	o := options.Common{}
	o.EnvConfig = readEnvConfig()

	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         "Emporous Client",
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

	cmd.AddCommand(NewInspectCmd(&o))
	cmd.AddCommand(NewBuildCmd(&o))
	cmd.AddCommand(NewPushCmd(&o))
	cmd.AddCommand(NewPullCmd(&o))
	cmd.AddCommand(NewServeCmd(&o))
	cmd.AddCommand(NewVersionCmd(&o))
	cmd.AddCommand(NewCreateCmd(&o))

	return cmd
}

func readEnvConfig() options.EnvConfig {
	envConfig := options.EnvConfig{}

	devModeString := os.Getenv("EMPOROUS_DEV_MODE")
	devMode, err := strconv.ParseBool(devModeString)
	envConfig.EMPOROUS_DEV_MODE = err == nil && devMode

	return envConfig
}
