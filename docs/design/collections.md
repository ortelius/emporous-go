Design: Collections
===
- [Design: Collections](#design-collections)
  - [Assumptions](#assumptions)
  - [Possibilities](#possibilities)


## Assumptions
- A collection is a node of nodes. This would be a directed graph with cycle detection.
- A collection represents one OCI artifact that can contain one to many descriptors that reference each other.
- Collections are built with the builder. They may have a "dependency" annotation is symbolize a connection to another collection or node type.


## Possibilities
- Collection and dependency collection are full resolved on the initial run of `client run`. Each collection will be compressed and stored in a cache and processed in batches.