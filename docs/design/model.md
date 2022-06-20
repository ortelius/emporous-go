Design: Model
===
- [Design: Model](#design-model)
- [Interface types defined in the model package.](#interface-types-defined-in-the-model-package)
  - [Traversal](#traversal)


The model package and the sub-package contains all types and methods that can be used to define and work with UOR data.

Package and sub-package layout 
```
model/
├── doc.go
├── nodes
│   ├── basic
│   │   └── basic.go
│   ├── collection
│   │   ├── assembly.go
│   │   ├── builder.go
│   │   ├── collection.go
│   │   ├── collection_test.go
│   │   ├── doc.go
│   │   ├── edge.go
│   │   ├── iterator.go
│   │   └── iterator_test.go
│   └── doc.go
├── traversal
│   ├── budget.go
│   ├── errors.go
│   ├── traversal.go
│   └── traversal_test.go
└── types.go
```

> For more information on the concrete node types, see [nodes](nodes.md).

# Interface types defined in the model package.
- Tree: This structure defines relationship between different UOR node types (i.e. connects different UOR collections). A UOR collection can reference one to many nodes as a dependency.
- Node: This represents a single addressable unit. This can contain one to many files. This interface is a read-only contains read only methods.
- NodeBuilder: This container methods for building immutable nodes.
- Iterator: Nodes can be iterable (e.g. a UOR collection). Using the iterator interface allows these structures to be iterated over during tree traversal.
- Matcher: Defines criteria for node searching in a tree or subtree.
- Attributes: This represents the methods that would be used by a structure containing a set of attributes.

## Traversal

The traversal package implements basic DFS and BFS traversal on model.Tree types. If a model.Node types is iterable (i.e. implements the mode.Iterator interface), the nodes within the types will be iterated over as well. More advance traversal methods will be coming in later releases.

Nodes can be skipped by returning the `ErrSkip` error in the `VisitFunc`. 


