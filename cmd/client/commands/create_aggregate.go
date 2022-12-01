package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	specgo "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2/registry"

	"github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// AggregateOptions describe configuration options that can
// be set using the push subcommand.
type AggregateOptions struct {
	*CreateOptions
	AttributeQuery string
	RegistryHost   string
	SchemaID       string
}

var clientAggregateExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Build aggregate from a query."},
	CommandString: "aggregate localhost:5001 myquery.yaml",
}

// NewAggregateCmd creates a new cobra.Command for the aggregate subcommand.
func NewAggregateCmd(createOps *CreateOptions) *cobra.Command {
	o := AggregateOptions{CreateOptions: createOps}

	cmd := &cobra.Command{
		Use:           "aggregate HOST QUERY",
		Short:         "Create an artifact aggregate from an attribute query",
		Example:       examples.FormatExamples(clientAggregateExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVarP(&o.SchemaID, "schema-id", "s", schema.UnknownSchemaID, "Schema ID to scope attribute query. Default is \"unknown\"")

	return cmd
}

func (o *AggregateOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RegistryHost = args[0]
	o.AttributeQuery = args[1]
	return nil
}

func (o *AggregateOptions) Validate() error {
	return nil
}

func (o *AggregateOptions) Run(ctx context.Context) error {
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	userQuery, err := config.ReadAttributeQuery(o.AttributeQuery)
	if err != nil {
		return err
	}

	var result ocispec.Index
	var queryJSON []byte
	if len(userQuery.Attributes) != 0 {
		o.Logger.Infof("Running attribute query")
		// Based on the way descriptor are written the schema is always the root key.
		constructedQuery := map[string]v1alpha1.Attributes{
			o.SchemaID: userQuery.Attributes,
		}
		queryJSON, err = json.Marshal(constructedQuery)
		if err != nil {
			return err
		}
		result, err = client.ResolveQuery(ctx, o.RegistryHost, nil, nil, queryJSON)
		if err != nil {
			return err
		}
		resultsJSON, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	if len(userQuery.Digests) != 0 {
		o.Logger.Infof("Running digest query")
		result, err = client.ResolveQuery(ctx, o.RegistryHost, nil, userQuery.Digests, nil)
		if err != nil {
			return err
		}
		resultsJSON, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	matcher := matchers.PartialAttributeMatcher{}
	if userQuery.LinkQuery.FilterBy != nil {
		attributeSet, err := config.ConvertToModel(userQuery.LinkQuery.FilterBy)
		if err != nil {
			return err
		}
		matcher = attributeSet.List()
	}

	if userQuery.LinkQuery.LinksTo != nil {
		o.Logger.Infof("Running links to query")
		result, err = client.ResolveQuery(ctx, o.RegistryHost, userQuery.LinkQuery.LinksTo, nil, nil)
		if err != nil {
			return err
		}

		var collections []collection.Collection
		for _, desc := range result.Manifests {
			if desc.Annotations == nil {
				continue
			}

			hint, ok := desc.Annotations["namespaceHint"]
			if !ok {
				continue
			}

			constructedRef := fmt.Sprintf("%s/%s@%s", o.RegistryHost, hint, desc.Digest)
			collection, err := client.LoadCollection(ctx, constructedRef)
			if err != nil {
				return err
			}
			collections = append(collections, collection)
		}

		filterResults, err := filterIndex(collections, matcher)
		if err != nil {
			return err
		}
		resultsJSON, err := json.MarshalIndent(filterResults, " ", " ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	if userQuery.LinkQuery.LinksFrom != nil {
		o.Logger.Infof("Running links from query")
		result, err = client.ResolveQuery(ctx, o.RegistryHost, nil, userQuery.LinkQuery.LinksFrom, nil)
		if err != nil {
			return err
		}

		var descs []ocispec.Descriptor
		for _, manifest := range result.Manifests {
			if manifest.Annotations != nil {
				link, ok := manifest.Annotations[uorspec.AnnotationLink]
				if ok {
					if err := json.Unmarshal([]byte(link), &descs); err != nil {
						return err
					}
					o.Logger.Infof("Found links for digest %s", manifest.Digest)
				}
			}

			var collections []collection.Collection
			for _, desc := range descs {
				node, err := v2.NewNode(desc.Digest.String(), desc)
				if err != nil {
					return err
				}
				if node.Properties != nil && node.Properties.IsALink() {
					constructedRef := fmt.Sprintf("%s/%s@%s", node.Properties.Link.RegistryHint, node.Properties.Link.NamespaceHint, desc.Digest)
					collection, err := client.LoadCollection(ctx, constructedRef)
					if err != nil {
						return err
					}
					collections = append(collections, collection)
				}
			}
			filterResults, err := filterIndex(collections, matcher)
			if err != nil {
				return err
			}
			resultsJSON, err := json.MarshalIndent(filterResults, " ", " ")
			if err != nil {
				return err
			}
			fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
		}
	}

	return nil
}

func filterIndex(collections []collection.Collection, matcher model.Matcher) (ocispec.Index, error) {
	filteredIndex := ocispec.Index{
		Versioned: specgo.Versioned{
			SchemaVersion: 2,
		},
		MediaType: ocispec.MediaTypeImageIndex,
	}
	for _, currCol := range collections {
		subCollection, err := currCol.SubCollection(matcher)
		if err != nil {
			return filteredIndex, err
		}
		if len(subCollection.Nodes()) != 0 {
			rootNode, err := currCol.Root()
			if err != nil {
				return filteredIndex, err
			}

			rootDesc, ok := rootNode.(*v2.Node)
			if ok {
				rootOCIDesc := rootDesc.Descriptor()
				ref, err := registry.ParseReference(currCol.Address())
				if err != nil {
					return filteredIndex, err
				}
				props := descriptor.Properties{
					Link: &uorspec.LinkAttributes{
						RegistryHint:  ref.Registry,
						NamespaceHint: ref.Repository,
					},
				}
				propsJSON, err := props.MarshalJSON()
				if err != nil {
					return filteredIndex, err
				}
				rootOCIDesc.Annotations = map[string]string{
					uorspec.AnnotationUORAttributes: string(propsJSON),
				}
				filteredIndex.Manifests = append(filteredIndex.Manifests, rootOCIDesc)
			}
		}
	}
	return filteredIndex, nil
}
