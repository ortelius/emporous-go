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
	// TODO(jpower432): check that the directory passed exists and
	// validate the reference is valid
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
			node.Template, node.Links, err = p.GetLinkableData(buf.Bytes(), fileIndex)
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
	client := registryclient.NewORASClient(o.Reference)
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
