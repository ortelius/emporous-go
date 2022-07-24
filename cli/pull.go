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
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/attributes"
	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/model/nodes/basic"
	"github.com/uor-framework/client/model/nodes/collection"
	"github.com/uor-framework/client/registryclient"
	"github.com/uor-framework/client/registryclient/orasclient"
	"github.com/uor-framework/client/util/workspace"
)

// PullOptions describe configuration options that can
// be set using the pull subcommand.
type PullOptions struct {
	*RootOptions
	Source     string
	Output     string
	Insecure   bool
	PlainHTTP  bool
	Configs    []string
	Attributes map[string]string
}

var clientPullExamples = templates.Examples(
	`
	# Push artifacts
	client pull localhost:5000/myartifacts:latest my-output-directory
	`,
)

// NewPullCmd creates a new cobra.Command for the pull subcommand.
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
	cmd.Flags().StringToStringVarP(&o.Attributes, "attributes", "", o.Attributes, "list of key,value pairs (e.g. key=value) for "+
		"retrieving artifacts by attributes")

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
	var pullFn pullFunc
	if o.Attributes != nil {
		pullFn = withAttributes
	} else {
		pullFn = func(ctx context.Context, po PullOptions) (ocispec.Descriptor, error) {
			desc, _, err := o.pullCollection(ctx, o.Output)
			return desc, err
		}
	}
	return o.run(ctx, pullFn)
}

func (o *PullOptions) run(ctx context.Context, pullFn pullFunc) error {
	desc, err := pullFn(ctx, *o)
	if err != nil {
		return err
	}
	o.Logger.Infof("Artifact %s from %s pulled to %s\n", desc.Digest, o.Source, o.Output)
	return nil
}

type pullFunc func(context.Context, PullOptions) (ocispec.Descriptor, error)

func withAttributes(ctx context.Context, o PullOptions) (ocispec.Descriptor, error) {
	cleanup, dir, err := o.mktempDir(o.Output)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	defer cleanup()

	desc, layerDescs, err := o.pullCollection(ctx, dir)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	// Convert descriptor annotations to type Attribute.
	attributesByFile := make(map[string]model.Attributes, len(layerDescs))
	for _, ldesc := range layerDescs {
		filename, ok := ldesc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			continue
		}
		attr := attributes.AnnotationsToAttributes(ldesc.Annotations)
		o.Logger.Debugf("Adding attributes %s for file %s", attr.String(), filename)
		attributesByFile[filename] = attr
	}

	space, err := workspace.NewLocalWorkspace(dir)
	if err != nil {
		return ocispec.Descriptor{}, err
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
		return ocispec.Descriptor{}, err
	}

	// Iterator through the collection using an iterator instead
	// of constructing a tree data structure for now since
	// we are handling one collection at first.
	itr := collection.NewByAttributesIterator(nodes)

	return desc, o.moveToResults(itr, o.Attributes)
}

// moveToResult will iterate through the collection
func (o *PullOptions) moveToResults(itr model.Iterator, matcher attributes.PartialAttributeMatcher) error {
	for itr.Next() {
		node := itr.Node()
		if matcher.Matches(node) {
			o.Logger.Debugf("Found match: %q", node.ID())
			newLoc := filepath.Join(o.Output, node.ID())
			cleanLoc := filepath.Clean(newLoc)
			if err := os.MkdirAll(filepath.Dir(cleanLoc), 0750); err != nil {
				return err
			}
			if err := os.Rename(node.Address(), cleanLoc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *PullOptions) pullCollection(ctx context.Context, output string) (ocispec.Descriptor, []ocispec.Descriptor, error) {
	var layerDescs []ocispec.Descriptor
	var mu sync.Mutex
	layerFn := func(_ context.Context, desc ocispec.Descriptor) error {
		mu.Lock()
		defer mu.Unlock()
		layerDescs = append(layerDescs, desc)
		return nil
	}
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithOutputDir(output),
		orasclient.WithPostCopy(layerFn),
	)
	if err != nil {
		return ocispec.Descriptor{}, nil, fmt.Errorf("error configuring client: %v", err)
	}

	desc, err := client.Execute(ctx, o.Source, registryclient.TypePull)
	if err != nil {
		return ocispec.Descriptor{}, nil, err
	}
	return desc, layerDescs, client.Destroy()
}

// mkTempDir will make a temporary dir and return the name
// and cleanup method
func (o *PullOptions) mktempDir(parent string) (func(), string, error) {
	dir, err := ioutil.TempDir(parent, "collection.*")
	return func() {
		if err := os.RemoveAll(dir); err != nil {
			o.Logger.Fatalf(err.Error())
		}
	}, dir, err
}
