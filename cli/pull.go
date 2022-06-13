package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/registryclient"
	"github.com/uor-framework/client/registryclient/orasclient"
)

type PullOptions struct {
	*RootOptions
	Source    string
	Output    string
	Insecure  bool
	PlainHTTP bool
	Configs   []string
	DSConfig  string
}

var clientPullExamples = templates.Examples(
	`
	# Push artifacts
	client pull localhost:5000/myartifacts:latest my-output-directory
	`,
)

func NewPullCmd(rootOpts *RootOptions) *cobra.Command {
	o := PullOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "pull SRC DST",
		Short:         "Pull a UOR collection based on content or attribute address",
		Example:       clientPullExamples,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringArrayVarP(&o.Configs, "auth-configs", "c", o.Configs, "auth config paths")
	cmd.Flags().BoolVarP(&o.Insecure, "insecure", "", o.Insecure, "allow connections to SSL registry without certs")
	cmd.Flags().BoolVarP(&o.PlainHTTP, "plain-http", "", o.PlainHTTP, "use plain http and not https")
	cmd.Flags().StringVarP(&o.DSConfig, "dsconfig", "", o.DSConfig, "DataSet config path")

	return cmd
}

func (o *PullOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.Source = args[0]
	o.Output = args[1]
	return nil
}

func (o *PullOptions) Validate() error {
	if _, err := os.Stat(o.Output); err != nil {
		if err := os.MkdirAll(o.Output, 0750); err != nil {
			return err
		}
	}
	return nil
}

func (o *PullOptions) Run(ctx context.Context) error {
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithOutputDir(o.Output),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}

	desc, err := client.Execute(ctx, o.Source, registryclient.TypePull)
	if err != nil {
		return err
	}

	o.Logger.Infof("Artifact %s from %s pulled to %s\n", desc.Digest, o.Source, o.Output)

	return nil
}
