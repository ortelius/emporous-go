Design: Node Types
===
- [Design: Node Types](#design-node-types)
- [Type Requirements](#type-requirements)
- [Basic Node](#basic-node)
- [Collection Node](#collection-node)

# Type Requirements

The model.Nodes types must implement the following methods:

- ID: This is a string ID that **MUST** be unique with a implemented model.Tree or model.Collection (if applicable).
- Address: This represent the location of the node whether it be a relative or absolute path on disk or the registry path for a remote location.
- Attributes: This is a collection of key, values pair that represent information about the node. Example for a file could be name, size, and file type.

# Basic Node

A basic node represent the smallest unit of information in a UOR dataset. This can be part of a tree alongside Collection node
types or references within a Collection node.

# Collection Node

A Collection node implements a node of nodes patterns. A collection node implements the model.Node interface itself so its relationship to other nodes can be described.

Quick information about the collection node:
- It represent one workspace locally and on OCI artifact remotely.
- It store a collection of nodes in a rooted graph data structure.
- As a node it presents the attributes of the root node, with the Attribute method is called.

The reference implementation use case for the Collection node is the following:

The Collection node represent a workspace with files that contained links to other files within the workspace. In order to render these
files in a way that would be consumable in a registry, the relationship between the files must be tracked. This is the way workspaces are created
in the builder.Compatibility mode.
