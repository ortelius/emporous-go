# Collection Workflows

<!--toc-->
- [Collection Workflows](#collection-workflows)
    * [Collection Publishing](#collection-publishing)
    * [Collection Pulling](#collection-pulling)

<!-- tocstop -->

## Collection Publishing

The workflow for collection publishing is very similar to the workflow used to build container images with most tooling options.
If the schema is published before the collection and associated to the collection, the specified attributes in the dataset configuration
will be validated against the schema during collection building.
This is demonstrated below by the diagram.

```mermaid
graph LR;
A[Build]-->B{Collection or Schema?};
B -- Schema --> C[Build Schema]
B -- Collection --> D[Build Collection]
C --> E[Push Schema]
E --> D
D --> F[Push Collection]
```

This would be the workflow if the ultimate goal was to publish a collection with or without a schema. Schema can be published
without a subsequent collection publish for later use as well.

```mermaid
graph LR;
A[Build]-->B{Collection or Schema?};
B -- Schema --> C[Build Schema]
B -- Collection --> D[Build Collection]
C --> E[Push]
D --> E
```

## Collection Pulling

Collections can be pulled as an entire OCI Artifact or filtered by an Attribute Query. The filtered OCI artifact is stored
in the build cache with the original manifest intact (sparse manifest) and the non-matching blobs (files) are not pulled into the cache.
All matching files are written to the cache and written to a user specified location.

The use of sparse manifest can pose a problem if re-tagging collections becomes part of the command line functionality in the future.
Some registries will reject manifests without all the blobs present. In this case, it may be of interest to reconstruct the manifest before pushing
and allow a flag to preserve the manifest, if desired.

## Collection Manager

Collection publishing and pulling can be accomplished over gRPC with the `serve` command. 
The client acts as a gRPC server that will retrieve and publish collections upon client request. 
A top-level type called `Manager` is used to provide this functionality to the CLI and gRPC server. There is a default implementation located 
in the `manager` package that is currently used.

The gRPC server is reading and writing in locations relative to its instantiated location. Due to this, a unix domain socket is used for client/server communication. gRPC client
must provide absolute pathing for expected results.