# `nodes`

This README outlines how the various node types are used within the Emporous Collection. Each implements the `model.Node` interface and are used to represent the smallest unit of information in a Emporous Collection.

The model.Node types must implement the following methods:

- ID: This is a string ID that **MUST** be unique with the implemented graph are sub-graph.
- Address: This represent the location of the node whether it be a relative or absolute path on disk or the registry path for a remote location.
- Attributes: This is a collection of key, values pair that represent information about the node. Example for a file could be name, size, and file type.

## Basic Node

A basic node represent a generic node implementation that can be used in a Collection to represent addressable nodes that are not OCI descriptors. This is useful for representing directories, symlinks, and other non-descriptor nodes.

## Descriptor Node

A descriptor node can be used to reference a OCI descriptor in a Collection. Conversions are available in the package to convert a Collection node to an OCI node and to convert attributes from the annotations to the model.AttributeSet.

Descriptors properties have core attributes define by the OCI spec and collection specific attributes that are used to determine node type and runtime information. The can also contain additional attributes that are grouped by schema ID.

Core attributes are:

- Runtime (OCI Spec ImageConfig)
- Link (Emporous LinkAttributes)
- Descriptor (Emporous DescriptorAttributes)
- Schema (Emporous SchemaAttributes)
- File (Emporous File)

A JSON based matcher implementation can be found here to identify nodes based on their attributes in JSON format.

Usage Example

```go
mockAttributes := attributes.Attributes{
		"kind":    attributes.NewString("kind", "jpg"),
		"name":    attributes.NewString("name", "fish.jpg"),
		"another": attributes.NewString("another", "attribute"),
	}

n := FakeNode{A: mockAttributes}
m := JSONSubsetMatcher(`{"name":"fish.jpg"}`)
match, err := m.Matches(n)
```

It is currently being used to match nodes based on attributes when inspecting the build cache and when resolving nodes from a remote registry.


## Collection Node

A collection node is a node that represents a Collection. This represent a group of references and in the case of an artifact, a single OCI artifact. The address is the location of the root manifest.

### Usage

Collection can be built in memory or loaded from an OCI artifact. The following sections outline how to create and load a collection.

#### Loading a Collection

```go

import (
    "context"
    "encoding/json"

    ocispec "github.com/opencontainers/image-spec/specs-go/v1"
    "oras.land/oras-go/v2/content"
    "oras.land/oras-go/v2/registry/remote"
)

// Using oras library to load a collection from a remote registry

reference := "registry.example.com/repo"

repo, err := remote.NewRepository(reference)
if err != nil {
    return err
}

rootDesc, rc , err := repo.FetchReference(ctx, reference)
if err != nil {
    return err
}
defer rc.Close()

// Example fetcher function using oras library to fetch a child descriptors
fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
    r, err := repo.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}
	return content.ReadAll(r, desc)
}

c := collection.New(reference)
if err := collectionloader.LoadFromManifest(ctx, c, fetcherFn, rootDesc); err != nil {
    return err
}
```

#### Basic Collection Traversal using the `model` package

The `model` package provides a basic traversal implementation that can be used to traverse a `model.Node` type. The following example shows how to traverse a collection and print the node ID.

```go

import "fmt"

// Load a collection (see above)
graph := collection.New(reference)

...

root, err := graph.Root()
if err != nil {
    return nil, err
}

tracker := traversal.NewTracker(root, nil)
handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {
    fmt.Println(node.ID())
    return graph.From(node.ID())
	})
if err := tracker.Walk(ctx, handler, root); err != nil {
		return err
}
```


