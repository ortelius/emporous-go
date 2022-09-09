package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"

	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// PullOptions describe configuration options that can
// be set using the pull subcommand.
type PullOptions struct {
	*RootOptions
	Source         string
	Output         string
	Insecure       bool
	PullAll        bool
	PlainHTTP      bool
	Configs        []string
	AttributeQuery string
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
func NewPullCmd(rootOpts *RootOptions) *cobra.Command {
	o := PullOptions{RootOptions: rootOpts}

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

	cmd.Flags().StringArrayVarP(&o.Configs, "configs", "c", o.Configs, "auth config paths when contacting registries")
	cmd.Flags().BoolVarP(&o.Insecure, "insecure", "i", o.Insecure, "allow connections to SSL registry without certs")
	cmd.Flags().BoolVar(&o.PlainHTTP, "plain-http", o.PlainHTTP, "use plain http and not https when contacting registries")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "output location for artifacts")
	cmd.Flags().StringVar(&o.AttributeQuery, "attributes", o.AttributeQuery, "attribute query config path")
	cmd.Flags().BoolVar(&o.PullAll, "pull-all", o.PullAll, "pull all linked collections")

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
	o.Logger.Infof("Resolving artifacts for reference %s", o.Source)
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

	manifestDescs, err := o.pullCollections(ctx, matcher)
	if err != nil {
		return err
	}

	for _, desc := range manifestDescs {
		o.Logger.Infof("Artifact %s pulled to %s\n", desc.Digest, o.Output)
	}

	return nil
}

// pullCollection pulls a single collection and returns the manifest descriptors and an error.
func (o *PullOptions) pullCollection(ctx context.Context, graph *collection.Collection, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	// Filter the collection per the matcher criteria
	if matcher != nil {
		var matchedLeaf int
		matchFn := model.MatcherFunc(func(node model.Node) (bool, error) {
			// This check ensure we are not weeding out any manifests needed
			// for OCI DAG traversal.
			if len(graph.From(node.ID())) != 0 {
				return true, nil
			}

			// Check that this is a descriptor node and the blob is
			// not a config resource.
			desc, ok := node.(*descriptor.Node)
			if !ok {
				return false, nil
			}
			mediaType := desc.Descriptor().MediaType
			if mediaType == ocimanifest.UORConfigMediaType || mediaType == ocispec.MediaTypeImageConfig {
				return true, nil
			}

			match, err := matcher.Matches(node)
			if err != nil {
				return false, err
			}

			if match {
				matchedLeaf++
			}

			return match, nil
		})

		var err error
		*graph, err = graph.SubCollection(matchFn)
		if err != nil {
			return nil, err
		}

		if matchedLeaf == 0 {
			o.Logger.Infof("No matches found for collection %s", graph.Address())
			return nil, nil
		}
	}

	var mu sync.Mutex
	successorFn := func(_ context.Context, fetcher orascontent.Fetcher, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		mu.Lock()
		successors := graph.From(desc.Digest.String())
		mu.Unlock()
		var result []ocispec.Descriptor
		for _, s := range successors {
			d, ok := s.(*descriptor.Node)
			if ok {
				result = append(result, d.Descriptor())
			}
		}
		return result, nil
	}

	cache, err := layout.NewWithContext(ctx, o.cacheDir)
	if err != nil {
		return nil, err
	}
	client, err := orasclient.NewClient(
		orasclient.WithSuccessorFn(successorFn),
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithCache(cache),
	)
	if err != nil {
		return nil, fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	desc, err := client.Pull(ctx, graph.Address(), file.New(o.Output))
	if err != nil {
		return nil, err
	}

	return []ocispec.Descriptor{desc}, cache.Tag(ctx, desc, graph.Address())
}

// pullCollections pulls one or more collections and returns the manifest descriptors and an error.
func (o *PullOptions) pullCollections(ctx context.Context, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
	)
	if err != nil {
		return nil, fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()
	root, err := o.loadFromReference(ctx, o.Source, client)
	if err != nil {
		return nil, err
	}
	return o.copy(ctx, root, client, matcher)
}

// copy performs graph traversal of linked collections and performs collection copies filtered by the matcher.
func (o *PullOptions) copy(ctx context.Context, root model.Node, client registryclient.Remote, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	seen := map[string]struct{}{}
	var manifestDesc []ocispec.Descriptor

	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {

		descs, err := o.pullCollection(ctx, node.(*collection.Collection), matcher)
		if err != nil {
			return nil, err
		}
		manifestDesc = append(manifestDesc, descs...)

		// Do not descend into the linked collection graph if
		// pull all is false.
		if !o.PullAll {
			return nil, nil
		}

		successors, err := o.getSuccessors(ctx, node.Address(), client)
		if err != nil {
			if errors.Is(err, ocimanifest.ErrNoCollectionLinks) {
				o.Logger.Debugf("collection %s has no links", node.Address())
				return nil, nil
			}
			return nil, err
		}

		var result []model.Node
		for _, s := range successors {
			if _, found := seen[s]; !found {
				o.Logger.Debugf("found link %s for collection %s", s, node.Address())
				childNode, err := o.loadFromReference(ctx, s, client)
				if err != nil {
					return nil, err
				}
				result = append(result, childNode)
				seen[s] = struct{}{}
			}
		}
		return result, nil
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return nil, err
	}

	return manifestDesc, nil
}

// getSuccessors retrieves all referenced collections from a source collection.
func (o *PullOptions) getSuccessors(ctx context.Context, reference string, client registryclient.Remote) ([]string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return nil, err
	}
	defer manBytes.Close()
	return ocimanifest.ResolveCollectionLinks(manBytes)
}

// loadFromReference loads a collection from an image reference.
func (o *PullOptions) loadFromReference(ctx context.Context, reference string, client registryclient.Remote) (*collection.Collection, error) {
	desc, _, err := client.GetManifest(ctx, reference)
	if err != nil {
		return nil, err
	}
	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return client.GetContent(ctx, reference, desc)
	}
	c := collection.New(reference)
	if err := collection.LoadFromManifest(ctx, c, fetcherFn, desc); err != nil {
		return nil, err
	}
	c.Location = reference
	return c, nil
}
