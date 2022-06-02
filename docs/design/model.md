Design: Model
===
- [Design: Model](#design-model)


The model package and the sub-package contains all types and methods that can be used to defines and work with UOR data.

- Tree: This structure defines relationship between different UOR node types (i.e connects different UOR collection). A UOR collection can reference one to many nodes as a dependency.
- Node: This represents a single addressable unit. This can contain one to many files. This interface is a read-only contains read only methods.
- NodeBuilder: This container methods for build immutable nodes.
- Iterator: Nodes can be iterable (e.g. a UOR collection). Using the iterator interface allows these structures to be iterated over during tree traversal.
- Matcher: Defines criteria for node searching in a tree or subtree.
- Attributes: The represent the method that would be used by a structure containing a set of attributes.


