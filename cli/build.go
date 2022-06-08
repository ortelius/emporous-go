package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/builder"
	"github.com/uor-framework/client/builder/api/v1alpha1"
	"github.com/uor-framework/client/builder/graph"
	"github.com/uor-framework/client/builder/parser"
	"github.com/uor-framework/client/util/workspace"
)

type BuildOptions struct {
	*RootOptions
	RootDir string
	Output  string
}

var clientBuildExamples = templates.Examples(
	`
	# Template content in a directory
	# The default workspace is "client-workspace" in the current working directory.
	client build my-directory

	# Template content into a specified output directory.
	client build my-directory --output my-workspace
	`,
)

func NewBuildCmd(rootOpts *RootOptions) *cobra.Command {
	o := BuildOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "build SRC",
		Short:         "Template and build files from a local directory into a UOR dataset",
		Example:       clientBuildExamples,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "location to stored templated files")

	return cmd
}

func (o *BuildOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.RootDir = args[0]
	if o.Output == "" {
		o.Output = "client-workspace"
	}
	return nil
}

func (o *BuildOptions) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}
	return nil
}

func (o *BuildOptions) Run(ctx context.Context) error {
	o.Logger.Infof("Using output directory %q", o.Output)
	userSpace, err := workspace.NewLocalWorkspace(o.RootDir)
	if err != nil {
		return err
	}

	g := graph.NewGraph()

	fileIndex := make(map[string]struct{})
	// Do the initial walk to get an index of what is in the workspace
	err = userSpace.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.IsDir() {
			return nil
		}

		fileIndex[path] = struct{}{}

		return nil
	})
	if err != nil {
		return err
	}

	// Function to determine whether the
	// data should be replaced in the template
	tFunc := func(value interface{}) bool {
		stringValue, ok := value.(string)
		if !ok {
			return false
		}

		// If the file is found in the workspace
		// return true
		_, found := fileIndex[stringValue]
		return found
	}

	for path := range fileIndex {
		o.Logger.Infof("Adding node %s\n", path)
		node := graph.NewNode(path)

		perr := &parser.ErrInvalidFormat{}
		buf := new(bytes.Buffer)
		if err := userSpace.ReadObject(ctx, path, buf); err != nil {
			return err
		}
		p, err := parser.ByContentType(path, buf.Bytes())
		switch {
		case err == nil:
			p.AddFuncs(tFunc)
			node.Template, node.Links, err = p.GetLinkableData(buf.Bytes())
			if err != nil {
				return err
			}
		case !errors.As(err, &perr):
			return err
		}

		g.Nodes[node.Name] = node
	}

	for _, node := range g.Nodes {
		for link, data := range node.Links {
			// Currently with the parsing implementation
			// all initial values are expected to represent
			// the file string data present in the content.
			// FIXME(jpower432): Making this assumption could lead
			// to bug when trying to translate links to a graph.
			fpath, ok := data.(string)
			if !ok {
				return fmt.Errorf("link %q: value should be of type string", link)
			}
			if err := g.AddEdge(node.Name, fpath); err != nil {
				return err
			}
		}
	}

	templateBuilder := builder.NewBuilder(userSpace)
	renderSpace, err := workspace.NewLocalWorkspace(o.Output)
	if err != nil {
		return err
	}

	if err := templateBuilder.Run(ctx, g, renderSpace); err != nil {
		return fmt.Errorf("error building content: %v", err)
	}

	_, _ = fmt.Fprintf(o.IOStreams.Out, "\nTo publish this content, run the following command:")
	_, _ = fmt.Fprintf(o.IOStreams.Out, "\nclient push %s IMAGE\n", o.Output)

	return nil
}

// AddDescriptors adds the attributes of each file listed in the config
// to the annotations of its respective descriptor.
func AddDescriptors(d []v1.Descriptor, c v1alpha1.DataSetConfiguration) ([]v1.Descriptor, error) {
	// For each descriptor
	for i1, desc := range d {
		// Get the filename of the block
		filename := desc.Annotations["org.opencontainers.image.title"]
		// For each file in the config
		for i2, file := range c.Files {
			// If the filename of the block matches the filename of the file in the config
			// If the config has a grouping declared, make a valid regex.
			if strings.Contains(file.File, "*") && !strings.Contains(file.File, ".*") {
				file.File = strings.Replace(file.File, "*", ".*", -1)
			}
			namesearch, err := regexp.Compile(file.File)
			if err != nil {
				return []v1.Descriptor{}, err
			}
			// Find the matching descriptor
			if namesearch.Match([]byte(filename)) {
				// Get the k/v pairs from the config and add them to the block's annotations.
				for k, v := range c.Files[i2].Attributes {
					d[i1].Annotations[k] = v
				}
			} else {
				// If the block does not have a corresponding config element, skip it.
				continue
			}
		}
	}
	return d, nil
}
