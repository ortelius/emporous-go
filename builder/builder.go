package builder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/opencontainers/go-digest"

	"github.com/uor-framework/client/builder/parser"
	"github.com/uor-framework/client/graph"
	"github.com/uor-framework/client/util/workspace"
)

// Builder defines methods for building UOR datasets.
type Builder interface {
	Run(context.Context, *graph.Graph, workspace.Workspace) error
}

// Compatibility renders and writes templates from the source workspace
// in compatibility mode.
// Compatibility mode renders the template for each node representing a
// file replacing any links with the content address of the linked file.
type Compatibility struct {
	// Source workspace to find files
	Source workspace.Workspace
	// Links association with the node ID
	Links map[string]map[string]interface{}
	// Templates association with the node ID
	Templates map[string]template.Template
}

var _ Builder = &Compatibility{}

// NewCompatibilityBuilder creates an new Builder from the source
// workspace.
func NewCompatibilityBuilder(source workspace.Workspace) Compatibility {
	return Compatibility{
		Source:    source,
		Links:     map[string]map[string]interface{}{},
		Templates: map[string]template.Template{},
	}
}

// Run traverses the graph with build nodes to render the file templates to the destination workspace.
// All nodes are expected have an underlying concrete type of BuildNode.
func (b Compatibility) Run(ctx context.Context, g *graph.Graph, destination workspace.Workspace) error {
	root, err := g.Root()
	if err != nil {
		return fmt.Errorf("error calculating root node: %v", err)
	}

	buildroot, ok := root.(*graph.BuildNode)
	if !ok {
		return errors.New("wrong node type for root node")
	}
	// Links store the calculated sub problem (i.e. link hashes)
	links := make(map[string]interface{})
	return b.makeTemplates(ctx, g, buildroot, destination, links)
}

// makeTemplates does recursive DFS traversal of the graph to generate digest values and template files.
func (b Compatibility) makeTemplates(ctx context.Context, g *graph.Graph, start *graph.BuildNode, destination workspace.Workspace, links map[string]interface{}) error {
	if start == nil {
		return nil
	}

	// Template and hash each child node to
	// calculate parent node information
	for _, n := range g.NodesFrom(start.ID()) {
		if _, found := links[n.ID()]; found {
			continue
		}

		buildnode, ok := n.(*graph.BuildNode)
		if !ok {
			return errors.New("wrong node type")
		}

		if err := b.makeTemplates(ctx, g, buildnode, destination, links); err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)
	if b.isBuildable(start.ID()) {
		// Update all links data with currently accumulated
		// digest values and render the new file from template.
		nodeLinks, ok := b.Links[start.ID()]
		if !ok {
			return fmt.Errorf("buildable node %s has no values", start.ID())
		}
		b.Links[start.ID()] = mergeLinkData(nodeLinks, links)
		if err := b.render(buf, start.ID()); err != nil {
			return err
		}
	} else {
		if err := b.Source.ReadObject(ctx, start.ID(), buf); err != nil {
			return err
		}
	}

	if err := destination.WriteObject(ctx, start.ID(), buf.Bytes()); err != nil {
		return err
	}

	// Must calculate the digest after writing the content of
	// the buffer to file because the FromReader method consumes the data.
	dgst, err := digest.FromReader(buf)
	if err != nil {
		return err
	}

	templateValue := parser.ConvertFilenameForGoTemplateValue(start.ID())
	links[templateValue] = dgst
	return nil
}

// isBuildable will check the node ID able to be built with
// the render method.
func (b *Compatibility) isBuildable(id string) bool {
	tmpl, ok := b.Templates[id]
	if !ok {
		return false
	}
	return tmpl != (template.Template{})
}

// render will render the template for the current state of the
// node at the specified ID.
func (b *Compatibility) render(w io.Writer, id string) error {
	template, ok := b.Templates[id]
	if !ok {
		return fmt.Errorf("no template associated with node %s", id)
	}

	values, ok := b.Links[id]
	if !ok {
		return fmt.Errorf("no links associated with node %s", id)
	}

	return template.Execute(w, values)
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
