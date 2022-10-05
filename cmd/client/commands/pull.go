package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/content/file"

	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/manager/defaultmanager"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// PullOptions describe configuration options that can
// be set using the pull subcommand.
type PullOptions struct {
	*options.Common
	options.Remote
	options.RemoteAuth
	Source         string
	Output         string
	PullAll        bool
	AttributeQuery string
	NoVerify       bool
}

var clientPullExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "pull localhost:5001/test:latest",
		Descriptions: []string{
			"Pull collection reference.",
		},
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "pull localhost:5001/test:latest --pull-all",
		Descriptions: []string{
			"Pull collection reference and all linked references.",
		},
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "pull localhost:5001/test:latest --attributes attribute-query.yaml",
		Descriptions: []string{
			"Pull all content from reference that satisfies the attribute query.",
		},
	},
}

// NewPullCmd creates a new cobra.Command for the pull subcommand.
func NewPullCmd(common *options.Common) *cobra.Command {
	o := PullOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "pull SRC",
		Short:         "Pull a UOR collection based on content or attribute address",
		Example:       examples.FormatExamples(clientPullExamples...),
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
	o.RemoteAuth.BindFlags(cmd.Flags())
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "output location for artifacts")
	cmd.Flags().StringVar(&o.AttributeQuery, "attributes", o.AttributeQuery, "attribute query config path")
	cmd.Flags().BoolVar(&o.PullAll, "pull-all", o.PullAll, "pull all linked collections")
	cmd.Flags().BoolVarP(&o.NoVerify, "no-verify", "", o.NoVerify, "skip collection signature verification")

	return cmd
}

func (o *PullOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.Source = args[0]
	if o.Output == "" {
		o.Output = "."
	}
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

	if !o.NoVerify {
		o.Logger.Infof("Checking signature of %s", o.Source)
		if err := verifyCollection(ctx, o.Source, o.RemoteAuth.Configs, o.Remote); err != nil {
			return err
		}

	}

	matcher := matchers.PartialAttributeMatcher{}
	if o.AttributeQuery != "" {
		query, err := config.ReadAttributeQuery(o.AttributeQuery)
		if err != nil {
			return err
		}

		attributeSet, err := config.ConvertToModel(query.Attributes)
		if err != nil {
			return err
		}
		matcher = attributeSet.List()
	}

	cache, err := layout.NewWithContext(ctx, o.CacheDir)
	if err != nil {
		return err
	}

	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithCache(cache),
		orasclient.WithPullableAttributes(matcher),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	manager := defaultmanager.New(cache, o.Logger)

	var digests []string
	if !o.PullAll {
		digests, err = manager.Pull(ctx, o.Source, client, file.New(o.Output))
	} else {
		digests, err = manager.PullAll(ctx, o.Source, client, file.New(o.Output))
	}
	if err != nil {
		return err
	}

	if len(digests) == 0 {
		o.Logger.Infof("No matching collections found for %s", o.Source)
		return nil
	}

	o.Logger.Infof("Copied collection(s) to %s", o.Output)

	return nil
}
