package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/pkg/content"

	"github.com/uor-framework/client/builder"
	"github.com/uor-framework/client/builder/graph"
	"github.com/uor-framework/client/builder/parser"
	"github.com/uor-framework/client/registryclient"
	"github.com/uor-framework/client/util/workspace"
)

type RootOptions struct {
	IOStreams genericclioptions.IOStreams
	Reference string
	RootDir   string
	Insecure  bool
	PlainHTTP bool
	Configs   []string
}

func NewRootCmd() *cobra.Command {
	o := RootOptions{}
	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	cmd := &cobra.Command{
		Use: fmt.Sprintf(
			"%s <directory> <reference>",
			filepath.Base(os.Args[0]),
		),
		Short:         "Templates, builds, and publishes OCI content",
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringArrayVarP(&o.Configs, "config", "c", nil, "auth config path")
	cmd.Flags().BoolVarP(&o.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	cmd.Flags().BoolVarP(&o.PlainHTTP, "plain-http", "", false, "use plain http and not https")

	return cmd
}

func (o *RootOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RootDir = args[0]
	o.Reference = args[1]
	return nil
}

func (o *RootOptions) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}

	// TODO(jpower432): Validate the reference
	return nil
}

func (o *RootOptions) Run(ctx context.Context) error {
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

		if info.Mode().IsDir() {
			return nil
		}

		fileIndex[path] = struct{}{}

		return nil
	})
	if err != nil {
		return err
	}

	// Function to determine whether the
	// data should be replace in the template
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
		_, _ = fmt.Fprintf(o.IOStreams.Out, "Adding node %s\n", path)
		node := graph.NewNode(path)

		perr := &parser.ErrInvalidFormat{}
		file := filepath.Base(path)
		p, err := parser.ByExtension(file)
		switch {
		case err == nil:
			buf := new(bytes.Buffer)
			if err := userSpace.ReadObject(ctx, path, buf); err != nil {
				return err
			}
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
			// to bug when trying to translate links to a graph. There
			// may also be a way to avoid this reflection.
			stringData, ok := data.(string)
			if !ok {
				return fmt.Errorf("link %q: value should be of type string", link)
			}
			if err := g.AddEdge(node.Name, stringData); err != nil {
				return err
			}
		}
	}

	// Create a temporary directory for rendering the template
	// content under root directory
	tmpdir := fmt.Sprintf("tmp.%d", time.Now().Unix())
	renderSpace, err := userSpace.NewDirectory(tmpdir)
	if err != nil {
		return err
	}
	// Clean up rendered directory
	defer func() {
		if err := userSpace.DeleteDirectory(tmpdir); err != nil {
			_, _ = fmt.Fprintln(o.IOStreams.ErrOut, err)
		}
	}()

	if err = builder.Build(ctx, g, userSpace, renderSpace); err != nil {
		return err
	}

	// Gather descriptors written to the temporary directory for publishing
	registryOpts := content.RegistryOptions{
		Insecure:  o.Insecure,
		PlainHTTP: o.PlainHTTP,
		Configs:   o.Configs,
	}
	client := registryclient.NewORASClient(o.Reference, nil, registryOpts)
	var files []string
	err = renderSpace.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.Mode().IsRegular() {
			p := renderSpace.Path(path)
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	descs, err := client.GatherDescriptors(files...)
	if err != nil {
		return err
	}

	configDesc, err := client.GenerateConfig(nil)
	if err != nil {
		return err
	}

	if err := client.GenerateManifest(configDesc, nil, descs...); err != nil {
		return err
	}

	return client.Execute(ctx)
}
