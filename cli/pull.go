package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/content/file"

	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/nodes/basic"
	"github.com/uor-framework/uor-client-go/model/nodes/collection"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
	"github.com/uor-framework/uor-client-go/util/workspace"
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
	var pullFn pullFunc
	if o.AttributeQuery != "" {
		pullFn = withAttributes
	} else {
		pullFn = func(ctx context.Context, _ PullOptions) ([]ocispec.Descriptor, error) {
			manifestDescs, _, err := o.pullCollection(ctx, o.Output)
			return manifestDescs, err
		}
	}
	return o.run(ctx, pullFn)
}

func (o *PullOptions) run(ctx context.Context, pullFn pullFunc) error {
	o.Logger.Infof("Resolving artifacts for reference %s", o.Source)
	manifestDescs, err := pullFn(ctx, *o)
	if err != nil {
		return err
	}

	for _, desc := range manifestDescs {
		o.Logger.Infof("Artifact %s pulled to %s\n", desc.Digest, o.Output)
	}

	return nil
}

type pullFunc func(context.Context, PullOptions) ([]ocispec.Descriptor, error)

func withAttributes(ctx context.Context, o PullOptions) ([]ocispec.Descriptor, error) {
	cleanup, dir, err := o.mktempDir(o.Output)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	manifestDescs, layerDescs, err := o.pullCollection(ctx, dir)
	if err != nil {
		return nil, err
	}

	// Convert descriptor annotations to type Attribute.
	attributesByFile := make(map[string]model.AttributeSet, len(layerDescs))
	for _, ldesc := range layerDescs {
		filename, ok := ldesc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			continue
		}
		attr, err := ocimanifest.AnnotationsToAttributeSet(ldesc.Annotations, nil)
		if err != nil {
			return nil, err
		}
		o.Logger.Debugf("Adding attributes %s for file %s", attr.AsJSON(), filename)
		attributesByFile[filename] = attr
	}

	space, err := workspace.NewLocalWorkspace(dir)
	if err != nil {
		return nil, err
	}

	var nodes []model.Node
	err = space.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.IsDir() {
			return nil
		}

		attr, ok := attributesByFile[path]
		if !ok {
			o.Logger.Debugf("No attributes found for file %s", path)
			return nil
		}

		node := basic.NewNode(path, attr)

		node.Location = space.Path(path)
		nodes = append(nodes, node)
		o.Logger.Debugf("Node %s added", path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	itr := collection.NewInOrderIterator(nodes)

	query, err := config.ReadAttributeQuery(o.AttributeQuery)
	if err != nil {
		return nil, err
	}

	attributeSet, err := config.ConvertToModel(query.Attributes)
	if err != nil {
		return nil, err
	}

	moved, err := o.moveToResults(itr, attributeSet.List())
	if err != nil {
		return nil, err
	}

	if moved == 0 {
		o.Logger.Infof("No matching artifacts found")
		// Do not return the manifest descriptors if not matches are found.
		// Anything pulled into the temporary directory will be deleted.
		return nil, nil
	}

	return manifestDescs, nil

}

// moveToResult will iterate through the collection and moved any nodes with matching artifacts to the
// output directory.
func (o *PullOptions) moveToResults(itr model.Iterator, matcher matchers.PartialAttributeMatcher) (total int, err error) {
	for itr.Next() {
		node := itr.Node()
		match, err := matcher.Matches(node)
		if err != nil {
			return total, err
		}
		if match {
			o.Logger.Debugf("Found match: %q", node.ID())
			newLoc := filepath.Join(o.Output, node.ID())
			cleanLoc := filepath.Clean(newLoc)
			if err := os.MkdirAll(filepath.Dir(cleanLoc), 0750); err != nil {
				return total, err
			}
			if err := os.Rename(node.Address(), cleanLoc); err != nil {
				return total, err
			}
			total++
		}
	}
	return total, nil
}

// pullCollection will pull one or more collections and return the manifest descriptors, layer descriptors, and an error.
func (o *PullOptions) pullCollection(ctx context.Context, output string) ([]ocispec.Descriptor, []ocispec.Descriptor, error) {
	var layerDescs []ocispec.Descriptor
	var manifestDescs []ocispec.Descriptor
	var mu sync.Mutex
	layerFn := func(_ context.Context, desc ocispec.Descriptor) error {
		mu.Lock()
		defer mu.Unlock()
		layerDescs = append(layerDescs, desc)
		return nil
	}

	cache, err := layout.NewWithContext(ctx, o.cacheDir)
	if err != nil {
		return manifestDescs, layerDescs, err
	}
	client, err := orasclient.NewClient(
		orasclient.WithPostCopy(layerFn),
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithCache(cache),
	)
	if err != nil {
		return manifestDescs, layerDescs, fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	pullSource := func(source string) (ocispec.Descriptor, error) {
		// TODO(jpower432): Write an method to pull blobs
		// by attribute from the cache to a content.Store.
		desc, err := client.Pull(ctx, source, file.New(output))
		if err != nil {
			return desc, fmt.Errorf("pull error for reference %s: %v", o.Source, err)
		}

		o.Logger.Debugf("Pulled down %s for reference %s", desc.Digest, source)

		// The cache will be populated by the pull command
		// Ensure the resource is captured in the index.json by
		// tagging the reference.
		if err := cache.Tag(ctx, desc, source); err != nil {
			return desc, err
		}
		return desc, nil
	}

	desc, err := pullSource(o.Source)
	if err != nil {
		return manifestDescs, layerDescs, err
	}
	manifestDescs = append(manifestDescs, desc)

	// Resolve source links and all linked collections by BFS.
	if o.PullAll {
		o.Logger.Infof("Resolving linked collections for reference %s", o.Source)
		visitedRefs := map[string]struct{}{o.Source: {}}
		linkedRefs, err := cache.ResolveLinks(ctx, o.Source)
		if err := o.checkResolvedLinksError(o.Source, err); err != nil {
			return manifestDescs, layerDescs, err
		}

		for len(linkedRefs) != 0 {
			currRef := linkedRefs[0]
			linkedRefs = linkedRefs[1:]
			if _, ok := visitedRefs[currRef]; ok {
				continue
			}
			visitedRefs[currRef] = struct{}{}

			desc, err := pullSource(currRef)
			if err != nil {
				return manifestDescs, layerDescs, err
			}
			manifestDescs = append(manifestDescs, desc)
			o.Logger.Infof("Resolving linked collections for reference %s", currRef)
			currLinks, err := cache.ResolveLinks(ctx, currRef)
			if err := o.checkResolvedLinksError(currRef, err); err != nil {
				return manifestDescs, layerDescs, err
			}
			linkedRefs = append(linkedRefs, currLinks...)
		}
	}

	return manifestDescs, layerDescs, nil
}

// checkResolvedLinksError logs errors when no collection is
// found.
func (o *PullOptions) checkResolvedLinksError(ref string, err error) error {
	if err == nil {
		return nil
	}
	if !errors.Is(err, ocimanifest.ErrNoCollectionLinks) {
		return err
	}
	o.Logger.Infof("No linked collections found for %s", ref)
	return nil
}

// mkTempDir will make a temporary dir and return the name
// and cleanup method.
func (o *PullOptions) mktempDir(parent string) (func(), string, error) {
	dir, err := ioutil.TempDir(parent, "collection.*")
	return func() {
		if err := os.RemoveAll(dir); err != nil {
			o.Logger.Fatalf(err.Error())
		}
	}, dir, err
}
