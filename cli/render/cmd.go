package render

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/builder"
	"github.com/uor-framework/uor-client-go/builder/parser"
	"github.com/uor-framework/uor-client-go/cli/options"
	"github.com/uor-framework/uor-client-go/model/nodes/basic"
	"github.com/uor-framework/uor-client-go/model/nodes/collection"
	"github.com/uor-framework/uor-client-go/util/examples"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// Options describe configuration options that can
// be set using the render subcommand.
type Options struct {
	*options.Common
	RootDir string
	Output  string
}

var clientRenderExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "render my-directory",
		Descriptions: []string{
			"Template content in a directory.",
			"The default workspace is \"client-workspace\" in the current working directory.",
		},
	},
	{
		Descriptions:  []string{"Template content into a specified output directory."},
		CommandString: "build my-directory --output my-workspace",
		RootCommand:   filepath.Base(os.Args[0]),
	},
}

// NewCmd creates a new cobra.Command for the render subcommand.
func NewCmd(commonOpts *options.Common) *cobra.Command {
	o := Options{Common: commonOpts}

	cmd := &cobra.Command{
		// TODO(sabre1041) Reenable/remove once build capability strategy determined
		Hidden:        !o.UOR_DEV_MODE,
		Use:           "render SRC",
		Short:         "Template and build files from a local directory into a UOR dataset",
		Example:       examples.FormatExamples(clientRenderExamples...),
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

func (o *Options) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.RootDir = args[0]
	if o.Output == "" {
		o.Output = "client-workspace"
	}
	return nil
}

func (o *Options) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}
	return nil
}

func (o *Options) Run(ctx context.Context) error {
	o.Logger.Infof("Using output directory %q", o.Output)
	userSpace, err := workspace.NewLocalWorkspace(o.RootDir)
	if err != nil {
		return err
	}

	c := collection.New(o.Output)

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

	templateBuilder := builder.NewCompatibilityBuilder(userSpace)

	for path := range fileIndex {
		o.Logger.Infof("Adding node %s\n", path)

		// Since the paths will be unique in this
		// case, the id is set as the location.
		node := basic.NewNode(path, nil)
		node.Location = path

		perr := &parser.ErrInvalidFormat{}
		buf := new(bytes.Buffer)
		if err := userSpace.ReadObject(ctx, path, buf); err != nil {
			return err
		}
		p, err := parser.ByContentType(path, buf.Bytes())
		switch {
		case err == nil:
			p.AddFuncs(tFunc)
			templates, links, err := p.GetLinkableData(buf.Bytes())
			if err != nil {
				return err
			}
			templateBuilder.Links[node.ID()] = links
			templateBuilder.Templates[node.ID()] = templates
		case !errors.As(err, &perr):
			return err
		}

		if err := c.AddNode(node); err != nil {
			return err
		}
	}

	for _, node := range c.Nodes() {
		for link, data := range templateBuilder.Links[node.ID()] {
			// Currently with the parsing implementation
			// all initial values are expected to represent
			// the file string data present in the content.
			// FIXME(jpower432): Making this assumption could lead
			// to bug when trying to translate links to a graph.
			fpath, ok := data.(string)
			if !ok {
				return fmt.Errorf("link %q: value should be of type string", link)
			}
			to := c.NodeByID(fpath)
			edge := collection.NewEdge(node, to)
			if err := c.AddEdge(edge); err != nil {
				return err
			}
		}
	}

	renderSpace, err := workspace.NewLocalWorkspace(o.Output)
	if err != nil {
		return err
	}
	if err := templateBuilder.Run(ctx, c, renderSpace); err != nil {
		return fmt.Errorf("error building content: %v", err)
	}

	_, _ = fmt.Fprintf(o.IOStreams.Out, "\nTo publish this content, run the following command:")
	_, _ = fmt.Fprintf(o.IOStreams.Out, "\nclient push %s IMAGE\n", o.Output)

	return nil
}
