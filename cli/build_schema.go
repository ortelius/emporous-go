package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/uor-framework/uor-client-go/schema"

	"github.com/spf13/cobra"

	load "github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// BuildSchemaOptions describe configuration options that can
// be set using the build schema subcommand.
type BuildSchemaOptions struct {
	*BuildOptions
	SchemaConfig string
}

var clientBuildSchemaExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build schema artifacts."},
		CommandString: "build schema schema-config.yaml localhost:5000/myartifacts:latest",
	},
}

// NewBuildSchemaCmd creates a new cobra.Command for the build schema subcommand.
func NewBuildSchemaCmd(buildOpts *BuildOptions) *cobra.Command {
	o := BuildSchemaOptions{BuildOptions: buildOpts}

	cmd := &cobra.Command{
		Use:           "schema CFG-PATH DST",
		Short:         "Build and save a UOR schema as an OCI artifact",
		Example:       examples.FormatExamples(clientBuildSchemaExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	return cmd
}

func (o *BuildSchemaOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.SchemaConfig = args[0]
	o.Destination = args[1]
	return nil
}

func (o *BuildSchemaOptions) Validate() error {
	info, err := os.Stat(o.SchemaConfig)
	if err != nil {
		return fmt.Errorf("schema configuration %q: %v", o.SchemaConfig, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("schema configuration %q: file is not regular", o.SchemaConfig)
	}
	return nil
}

func (o *BuildSchemaOptions) Run(ctx context.Context) error {

	config, err := load.ReadSchemaConfig(o.SchemaConfig)
	if err != nil {
		return err
	}

	cache, err := layout.NewWithContext(ctx, o.cacheDir)
	if err != nil {
		return err
	}

	client, err := orasclient.NewClient()
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}

	generatedSchema, err := schema.FromTypes(config.Schema.AttributeTypes)
	if err != nil {
		return err
	}

	desc, err := client.AddContent(ctx, ocimanifest.UORSchemaMediaType, generatedSchema.Export(), nil)
	if err != nil {
		return err
	}

	configDesc, err := client.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
	if err != nil {
		return err
	}

	_, err = client.AddManifest(ctx, o.Destination, configDesc, nil, desc)
	if err != nil {
		return err
	}

	desc, err = client.Save(ctx, o.Destination, cache)
	if err != nil {
		return fmt.Errorf("client save error for reference %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Schema %s built with reference name %s\n", desc.Digest, o.Destination)

	return client.Destroy()
}
