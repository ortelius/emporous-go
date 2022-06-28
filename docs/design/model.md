Design: Model
===
- [Design: Model](#design-model)
- [Interface types defined in the model package.](#interface-types-defined-in-the-model-package)


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
└── types.go
```

> For more information on the concrete node types, see [nodes](nodes.md).

# Interface types defined in the model package.
- Node: This represents a single addressable unit. This can contain one to many files. This interface is a read-only contains read only methods.
- NodeBuilder: This container methods for building immutable nodes.
- Iterator: Nodes can be iterable (e.g. a UOR collection). Using the iterator interface allows these structures to be iterated over during tree traversal.
- Matcher: Defines criteria for node searching in a tree or subtree.
- Attributes: This represents the methods that would be used by a structure containing a set of attributes.


