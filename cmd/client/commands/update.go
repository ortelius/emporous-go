package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/manager/defaultmanager"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// PushOptions describe configuration options that can
// be set using the push subcommand.

type UpdateOptions struct {
	PullOptions
	Add    string
	Delete string
}

var clientUpdateExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Update Collection."},
	CommandString: "update my-directory localhost:5000/myartifacts:latest",
}

// NewPushCmd creates a new cobra.Command for the push subcommand.
func NewUpdateCmd(common *options.Common) *cobra.Command {
	o := UpdateOptions{}

	o.PullOptions.Common = common

	cmd := &cobra.Command{
		Use:           "update SRC DST",
		Short:         "Update a UOR Collection",
		Example:       examples.FormatExamples(clientUpdateExamples),
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
	cmd.Flags().StringVarP(&o.Add, "add-content", "", o.Add, "add content to a UOR Collection")
	cmd.Flags().StringVarP(&o.AttributeQuery, "remove-content", "", o.AttributeQuery, "remove content from a UOR Collection")
	cmd.Flags().BoolVarP(&o.NoVerify, "no-verify", "", o.NoVerify, "skip collection signature verification")

	return cmd
}

func (o *UpdateOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.Source = args[0]
	return nil
}

func (o *UpdateOptions) Validate() error {
	return nil
}

func (o *UpdateOptions) Run(ctx context.Context) error {

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
		o.Logger.Infof("Attribute %s", attributeSet)

		matcher = attributeSet.List()
	}

	var space workspace.Workspace
	if o.Add != "" {
		var err error
		if space, err = workspace.NewLocalWorkspace(o.Add); err != nil {
			return err
		}
	}

	absCache, err := filepath.Abs(o.CacheDir)
	if err != nil {
		return err
	}
	cache, err := layout.NewWithContext(ctx, absCache)
	if err != nil {
		return err
	}

	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithPullableAttributes(matcher),
		orasclient.WithCache(cache),
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
	var remove bool
	if o.AttributeQuery != "" {
		fmt.Println("delete called")
		remove = true
	}
	var add bool
	if o.Add != "" {
		add = true
	}
	_, err = manager.Update(ctx, space, o.Source, add, remove, client)
	return err
}
