# Model

- [Model](#design-model)
- [Interface types defined in the model package.](#interface-types-defined-in-the-model-package)
- [Node](#node)
- [NodeBuilder](#nodebuilder)
- [Iterator](#iterator)
- [Matcher](#matcher)
- [Attribute](#attribute)
- [AttributeSet](#attributeset)

The model package and the sub-package contains all types and methods that can be used to define and work with Emporous data.

> For more information on the concrete node types, see [nodes](nodes.md).

# Interface types defined in the model package.

## Node

Node is an interface that is used to represent different types of self-describing addressable data (typically stored in
files locally or remotely). The methods defined in this interface are intended to be read-only. For methods that
manipulate or assemble nodes, see NodeBuilder.

## NodeBuilder

NodeBuilder is an interface that defines methods for building immutable Node types.

## Iterator

Nodes can be iterable (e.g. a Emporous Collection). Using the iterator interface allows these structures to be iterated over
during traversal.

## Matcher

Matcher is an interface that defines methods for node searching/matching that can guide Node graph traversal.

## Attribute

Attribute is an interface that defines a single attribute values with a key that is a type of string and a value that
can be a string, boolean, integer, number, or null value.

## AttributeSet

AttributeSet is an interface that represents the methods that would be used by a structure containing a set of Node
attributes.


