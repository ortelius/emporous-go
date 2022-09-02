package push

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cli/options"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// Options describe configuration options that can
// be set using the push subcommand.
type Options struct {
	*options.Common
	options.Remote
	Destination string
}

var clientPushExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Push artifacts."},
	CommandString: "push localhost:5000/myartifacts:latest",
}

// NewCmd creates a new cobra.Command for the push subcommand.
func NewCmd(commonOpts *options.Common) *cobra.Command {
	o := Options{Common: commonOpts}

	cmd := &cobra.Command{
		Use:           "push DST",
		Short:         "Push a UOR collection into a registry",
		Example:       examples.FormatExamples(clientPushExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	o.Remote.BindFlags(cmd.Flags())

	return cmd
}

func (o *Options) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.Destination = args[0]
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run(ctx context.Context) error {
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithAuthConfigs(o.Configs),
	)
	if err != nil {
		return err
	}

	cache, err := layout.NewWithContext(ctx, o.CacheDir)
	if err != nil {
		return err
	}

	desc, err := client.Push(ctx, cache, o.Destination)
	if err != nil {
		return fmt.Errorf("error publishing content to %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Artifact %s published to %s\n", desc.Digest, o.Destination)

	return client.Destroy()
}
