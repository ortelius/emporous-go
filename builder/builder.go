package builder

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/opencontainers/go-digest"

	"github.com/uor-framework/client/builder/graph"
	"github.com/uor-framework/client/builder/parser"
	"github.com/uor-framework/client/util/workspace"
)

// Builder renders and writes templates from the source workspace.
type Builder struct {
	Source workspace.Workspace
}

// NewBuilder creates an new Builder from the source
// workspace
func NewBuilder(source workspace.Workspace) Builder {
	return Builder{source}
}

// Run traverses the graph to render the file templates to the destination workspace.
func (b Builder) Run(ctx context.Context, g *graph.Graph, destination workspace.Workspace) error {
	root, err := g.Root()
	if err != nil {
		return fmt.Errorf("error calculating root node: %v", err)
	}
	// Links store the calculated sub problem (i.e. link hashes)
	links := make(map[string]interface{})
	return b.makeTemplates(ctx, g, root, destination, links)
}

// makeTemplates does recursive DFS traversal of the graph to generate digest values and template files.
func (b Builder) makeTemplates(ctx context.Context, g *graph.Graph, start *graph.Node, destination workspace.Workspace, links map[string]interface{}) error {
	if start == nil {
		return nil
	}

	// Template and hash each child node to
	// calculate parent node information
	for _, n := range start.Nodes {
		if _, found := links[n.Name]; found {
			continue
		}
		if err := b.makeTemplates(ctx, g, n, destination, links); err != nil {
			return err
		}
	}

	start.Links = mergeLinkData(start.Links, links)
	buf := new(bytes.Buffer)

	if start.Template != (template.Template{}) {
		if err := start.Template.Execute(buf, start.Links); err != nil {
			return err
		}
	} else {
		if err := b.Source.ReadObject(ctx, start.Name, buf); err != nil {
			return err
		}
	}

	if err := destination.WriteObject(ctx, start.Name, buf.Bytes()); err != nil {
		return err
	}

	// Must calculate the digest after writing the content of
	// the buffer to file because the FromReader method consumes the data.
	dgst, err := digest.FromReader(buf)
	if err != nil {
		return err
	}

	templateValue := parser.ConvertFilenameForGoTemplateValue(start.Name)
	links[templateValue] = dgst
	return nil
}

// mergeLinkData will merge any references to in-content links with
// the currently calculated values.
func mergeLinkData(in, curr map[string]interface{}) map[string]interface{} {
	for key := range in {
		currentVal, ok := curr[key]
		if ok {
			in[key] = currentVal
		}
	}
	return in
}
