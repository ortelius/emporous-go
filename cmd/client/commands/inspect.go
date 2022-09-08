package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/util/examples"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/content/layout"
)

// InspectOptions describe configuration options that can
// be set when using the inspect subcommand.
type InspectOptions struct {
	*options.Common
	Source         string
	AttributeQuery string
}

var clientInspectExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "inspect",
		Descriptions: []string{
			"List all references",
		},
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "inspect --reference localhost:5001/test:latest",
		Descriptions: []string{
			"List all descriptors for reference",
		},
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "inspect --reference localhost:5001/test:latest --attributes attribute-query.yaml",
		Descriptions: []string{
			"List all descriptors for reference with attribute filtering",
		},
	},
}

// NewInspectCmd creates a new cobra.Command for the inspect subcommand.
func NewInspectCmd(common *options.Common) *cobra.Command {
	o := InspectOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "inspect",
		Short:         "Print UOR collection information",
		Example:       examples.FormatExamples(clientInspectExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVarP(&o.AttributeQuery, "attributes", "a", o.AttributeQuery, "attribute query config path")
	cmd.Flags().StringVarP(&o.Source, "reference", "r", o.Source, "a reference to list descriptors for")

	return cmd
}

func (o *InspectOptions) Complete(_ []string) error {
	return nil
}

func (o *InspectOptions) Validate() error {
	if o.AttributeQuery != "" && o.Source == "" {
		return fmt.Errorf("must specify a reference with --reference")
	}
	return nil
}

func (o *InspectOptions) Run(ctx context.Context) error {
	cache, err := layout.NewWithContext(ctx, o.CacheDir)
	if err != nil {
		return err
	}

	if o.Source == "" {
		idx, err := cache.Index()
		if err != nil {
			return err
		}
		return o.formatManifestDescriptors(o.IOStreams.Out, idx.Manifests)
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

		o.Logger.Debugf("Resolving source %s to descriptor with %d attributes", o.Source, attributeSet.Len())

		matcher = attributeSet.List()
	}

	descs, err := cache.ResolveByAttribute(ctx, o.Source, matcher)
	if err != nil {
		return err
	}

	return o.formatDescriptors(o.IOStreams.Out, descs)
}

func (o *InspectOptions) formatManifestDescriptors(w io.Writer, descs []ocispec.Descriptor) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintf(tw, "Listing all references:\t%s\n", o.Source); err != nil {
		return err
	}
	for _, desc := range descs {
		if _, err := fmt.Fprintf(tw, "%s\n", desc.Annotations[ocispec.AnnotationRefName]); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (o *InspectOptions) formatDescriptors(w io.Writer, descs []ocispec.Descriptor) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintf(tw, "Listing matching descriptors for source:\t%s\n", o.Source); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(tw, "Name\tDigest\tSize\tMediaType"); err != nil {
		return err
	}
	for _, desc := range descs {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", desc.Annotations[ocispec.AnnotationTitle], desc.Digest, desc.Size, desc.MediaType); err != nil {
			return err
		}
	}
	return tw.Flush()
}
