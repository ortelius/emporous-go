package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/cmd/client/commands/options"
	load "github.com/emporous/emporous-go/config"
	"github.com/emporous/emporous-go/content/layout"
	"github.com/emporous/emporous-go/manager/defaultmanager"
	"github.com/emporous/emporous-go/registryclient/orasclient"
	"github.com/emporous/emporous-go/util/examples"
	"github.com/emporous/emporous-go/util/workspace"
)

// BuildCollectionOptions describe configuration options that can
// be set using the build collection subcommand.
type BuildCollectionOptions struct {
	*BuildOptions
	options.Remote
	options.RemoteAuth
	NoVerify bool
	RootDir  string
	// Dataset Config
	DSConfig string
}

var clientBuildCollectionExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build artifacts."},
		CommandString: "build collection my-directory localhost:5000/myartifacts:latest",
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build artifacts with custom annotations."},
		CommandString: "build collection my-directory localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml",
	},
}

// NewBuildCollectionCmd creates a new cobra.Command for the build collection subcommand.
func NewBuildCollectionCmd(buildOpts *BuildOptions) *cobra.Command {
	o := BuildCollectionOptions{BuildOptions: buildOpts}

	cmd := &cobra.Command{
		Use:           "collection SRC DST",
		Short:         "Build and save an OCI artifact from files",
		Example:       examples.FormatExamples(clientBuildCollectionExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	o.Remote.BindFlags(cmd.Flags())
	o.RemoteAuth.BindFlags(cmd.Flags())

	cmd.Flags().StringVarP(&o.DSConfig, "dsconfig", "d", o.DSConfig, "config path for artifact building and dataset configuration")
	cmd.Flags().BoolVar(&o.NoVerify, "no-verify", o.NoVerify, "skip schema signature verification")

	return cmd
}

func (o *BuildCollectionOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RootDir = args[0]
	o.Destination = args[1]
	return nil
}

func (o *BuildCollectionOptions) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}
	return nil
}

func (o *BuildCollectionOptions) Run(ctx context.Context) error {
	space, err := workspace.NewLocalWorkspace(o.RootDir)
	if err != nil {
		return err
	}

	absCache, err := filepath.Abs(o.CacheDir)
	if err != nil {
		return err
	}
	cache, err := layout.NewWithContext(ctx, absCache)
	if err != nil {
		return err
	}

	var clientOpts = []orasclient.ClientOption{
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
	}

	if !o.NoVerify {
		verificationFn := func(ctx context.Context, reference string) error {
			o.Logger.Debugf("Checking signature of %s", reference)
			err = verifyCollection(ctx, reference, o.RemoteAuth.Configs, o.Remote)
			if err != nil {
				return fmt.Errorf("collection %q: %v", reference, err)
			}
			return nil
		}
		clientOpts = append(clientOpts, orasclient.WithPrePullFunc(verificationFn))
	}

	client, err := orasclient.NewClient(clientOpts...)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	var config v1alpha1.DataSetConfiguration
	if len(o.DSConfig) > 0 {
		config, err = load.ReadDataSetConfig(o.DSConfig)
		if err != nil {
			return err
		}
	}

	manager := defaultmanager.New(cache, o.Logger)

	_, err = manager.Build(ctx, space, config, o.Destination, client)
	return err
}
