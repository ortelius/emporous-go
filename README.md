# Emporous

Note: The UR Client is being actively developed. Please consider joining the emporous Community to participate!

## Participate

Please join us in the discussion forum and feel free to ask questions about the emporous-Framework or UR Client.

## About

The Universal Runtime Client interacts with emporous artifacts and is aware of the runtime instruction
embedded in emporous artifacts.

To learn more about Universal Runtime visit the emporous Framework website at <https://universalreference.io>.

> WARNING: The repository is under active development and the API is subject to change.

## Development

### Requirements

- `go` version 1.18+

### Build

```
make
./bin/emporous -h
```

### Test

#### Unit

```
make test-unit
```

## Basic Usage

### Version

emporous version

### User Workflow

1. Use `emporous build schema` to build a schema to be used with a collection.
2. Use the `emporous build collection` command to build the workspace as an OCI artifact in build-cache The default location is ~/.emporous/cache. It can be set with the `emporous_CACHE` environment variable`.
3. Use the `emporous push` command to publish to a registry as an OCI artifact.
4. Use the `emporous pull` command to pull the artifact back to a local workspace.
5. Use the `emporous inspect` command to inspect the build cache to list information about references.

### Build a schema into an artifact

A schema can be optionally created prior to building a collection. Collections can then reference an already built schema or no schema at all.

```shell
emporous build schema schema-config.yaml localhost:5000/myschema:latest
```

### Build workspace into an artifact

Execute the following command to build a workspace into an an artifact:

```shell
emporous build collection my-workspace localhost:5000/myartifacts:latest
```

An optional dataset configuration file can be provided using the `--dsconfig` flag:

```shell
emporous build my-workspace localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml
```

### Push workspace to a registry location

Push a workspace to a remote registry

```shell
emporous push my-workspace localhost:5000/myartifacts:latest
```

### Pull emporous collection to a location

Pull a collection from a remote registry:

```shell
emporous pull localhost:5000/myartifacts:latest -o my-output-directory
```

### Pull subsets of a emporous collection to a location by attribute

Pull a portion of a collection by filtering for a set of attribute:

```shell
emporous pull localhost:5000/myartifacts:latest -o my-output-directory --attributes attribute-query.yaml
```

## Getting Started

This guide will walk through several exercises illustrating the use of the emporous Client

### Environment Setup

While these exercises can be implemented in any operating environment, this guide will make use of a prescriptive environment for which most users should be able to replicate easily.

1. Setting up an image registry

Most of the exercises require the interaction with a image registry. While there are several options available, [registry:2](https://hub.docker.com/_/registry) is the most straightforward to setup and use.

Using your container runtime of choice, start an instance of _registry:2_ (`docker` is used as the runtime here)

```shell
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

2. Set `EMPOROUS_CLIENT_GO_REPO` environment variable

Set an environment variable called `EMPOROUS_CLIENT_GO_REPO` with the current working directory within the repository. It will be used when working through the scenarios below.

```shell
export EMPOROUS_CLIENT_GO_REPO=$(pwd)
```

Create a directory called `examples` and change into this directory to begin working on the exercises.

```shell
mkdir ${EMPOROUS_CLIENT_GO_REPO}/exercises
cd ${EMPOROUS_CLIENT_GO_REPO}/exercises
```

3. Formatting json output using jq

[jq](https://stedolan.github.io/jq/) is a helpful utility for managing json content. It will be used to help format the output from requests to the remote registry.

Download and install the tool as described in the documentation.

### Basic Collection Publishing

The first exercise creates and publishes a basic collection to a remote repository

1. Create a new directory for this exercise called `basic` and change into this directory

```shell
mkdir -p ${EMPOROUS_CLIENT_GO_REPO}/exercises/basic
cd ${EMPOROUS_CLIENT_GO_REPO}/exercises/basic
```

2. Create a new _workspace_ directory for the collection called `basic-collection` and change into this directory

```shell
mkdir basic-collection
cd basic-collection
```

3. Add the content to the workspace. For example, two images of fishes an a text file.

Add the first photo to the current directory

```shell
cp ${EMPOROUS_CLIENT_GO_REPO}/test/fish.jpg .
```

Next, create a directory called `subdir1` to contain several additional files:

```shell
mkdir subdir1
```

Now, add a text file and another photo:

```shell
cp ${EMPOROUS_CLIENT_GO_REPO}/test/level1/fish2.jpg subdir1/
cp ${EMPOROUS_CLIENT_GO_REPO}/test/level1/file.txt subdir1/
```

4. Create a json document where the value of each kv pair is the path to each file within the directory. Multiple json documents can be used to create deep graphs, but a graph must only have one root. Multiple json docs in a build directory is for advanced use cases which will not be covered in this basic example and most just need one json document.

Create a json document called `basic.json` in the current directory:

```bash
cat << EOF > basic.json
{
    "fish": "fish.jpg",
    "text": "subdir1/file.txt",
    "fish2": "subdir1/fish2.jpg"
}
EOF
```

5. A Dataset Configuration can be use to assign _attributes_ to the various resources within a collection. This file must be located outside of the content directory and refer to the relative paths within the content directory. Add user defined key value pairs as subkeys to the `annotations`section. Each file should have as many attributes as possible. Multiple files can be referenced by using the `*` wildcard.

Navigate up one directory and create a file called `dataset-config.yaml` to contain the Dataset Configuration for the collection:

```shell
cd ..
```

```bash
cat << EOF > dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  files:
    - file: "fish.jpg"
      attributes:
        animal: "fish"
        habitat: "ocean"
        size: "small"
        color: "blue"
    - file: "subdir1/file.txt"
      attributes:
        fiction: true  
        genre: "science fiction"
    - file: "*.jpg"
      attributes:
        custom: "customval"
EOF
```

6. Run the emporous client _build_ command referencing the dataset config, the content directory, and the destination registry location to the local cache. Each of the examples in this document will make use of the registry located at `localhost:5000` that was started as part [Environment Setup](#environment-setup) section.

```shell
emporous build collection basic-collection localhost:5000/exercises/basic:latest --dsconfig dataset-config.yaml 
```

7. Run the emporous _push_ command to publish the collection to the remote repository.

NOTE: Since the registry that was used does not exposed a secure transport method (HTTPS), the `--plain-http` flag will need to be specified whenever there is any interaction with the remote registry. Feel free to adjust accordingly to the remote registry that is being used.

```
emporous push --plain-http localhost:5000/exercises/basic:latest
```

8. Inspect the OCI manifest of the published collection. The `jq` tool can be used to format the response to make it more readable. Once again, if the remote registry is exposed using HTTP or non trusted certificates, adjust the curl command below accordingly:

```shell
curl -s -H "Accept: application/vnd.oci.image.manifest.v1+json" http://localhost:5000/v2/exercises/basic/manifests/latest | jq -r
```

Notice that each of the files in the workspace are represented as _Layers_ within the Manifest. In addition, the relative location within the workspace along with the _attributes_ are added as _annotations_.

```json
...
  "layers": [
    {
      "mediaType": "text/plain; charset=utf-8",
      "digest": "sha256:8b8843c2c23a94efafa834c7b52547aa2cba63ed517c7891eba5b7386330482b",
      "size": 4,
      "annotations": {
        "org.opencontainers.image.title": "subdir1/file.txt",
        "emporous.attributes": "{\"converted\":{\"org.opencontainers.image.title\":\"subdir1/file.txt\"}, \"unknown\"{\"fiction\":true,\"genre\":\"science fiction\"}}"
      }
...
```

> TIP 1: The two root level keys "unknown" and "converted are important. Unknown is the schema ID used when no schema is linked to the collection and converted is used when the attributes are converted from annotations. You will need this information later when completing queries.

> TIP 2: Some other significant schema IDs to know are the following: core-link, core-descriptor, core-runtime, core-schema, and core-file. Schemas are below.


9. The emporous _inspect_ subcommand can be used to view the contents of the local cache. By default, the cache is located at `~/.emporous/cache/`.

```shell
emporous inspect
```

Notice the collection built previously is now present within the cache

```
Listing all references:  
localhost:5000/exercises/basic:latest
```

10. The collection can be pulled from the remote registry to verify the content. Use the emporous _pull_ subcommand to a directory called _my-output-directory_ using the `-o` flag:

```shell
emporous pull --plain-http localhost:5000/exercises/basic:latest -o my-output-directory
```

11. Instead of retrieving an entire collection, a subset can be retrieved by creating a `AttributeQuery` resource.

Create a file called `attribute-query.yaml` in the current directory:

```bash
cat << EOF > attribute-query.yaml
kind: AttributeQuery
apiVersion: client.emporous-framework.io/v1alpha1
attributes:
  unknown:
    fiction: true
EOF
```

Now use the emporous _pull_ subcommand along with the `--attributes` option placing the contents in a directory titled `my-filtered-output-directory`:

```shell
emporous pull localhost:5000/exercises/basic:latest --plain-http -o my-filtered-output-directory --attributes attribute-query.yaml
```

Since the `fiction=true` attribute was associated with only the _file.txt_ file it was the only resource retrieved from the collection.

```shell
tree my-filtered-output-directory

my-filtered-output-directory
└── subdir1
    └── file.txt

1 directory, 1 file
```

### Collection Publishing with Schema

A _Schema_ can be used to define the attributes associated with a collection along with linking multiple collections.

Be sure that the `EMPOROUS_CLIENT_GO_REPO` environment variable is defined as described at the beginning of the exercises along with the `examples` directory.

1. Create a new directory called `schema` within the _exercises_ directory and change into this directory

```shell
mkdir -p ${EMPOROUS_CLIENT_GO_REPO}/exercises/schema
cd ${EMPOROUS_CLIENT_GO_REPO}/exercises/schema
```

2. Create the Schema Configuration in a file called `schema-config.yaml` to define attribute keys and types for corresponding collections:

```bash
cat << EOF > animalschema.json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "animal": {
      "type": "string"
    },
    "size": {
      "type": "string"
    },
    "color": {
      "type": "string"
    },
    "habitat": {
      "type": "string"
    },
    "mammal": {
      "type": "boolean"
    }
  },
  "required": [
    "animal",
    "size",
    "color",
    "habitat",
    "mammal"
  ]
}
EOF
```
```bash
cat << EOF > schema-config.yaml
kind: SchemaConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
schema:
  id: myschemaid
  schemaPath: animalschema.json
EOF
```

3. Use the emporous _build_ subcommand to build and save the schema within the local cache:

```shell
emporous build schema schema-config.yaml localhost:5000/exercises/myschema:latest
```

4. Push the schema to the remote registry:

```
emporous push --plain-http localhost:5000/exercises/myschema:latest
```

5. Create a new directory called `schema-collection` to contain a workspace to demonstrate how a schema can be used

```shell
mkdir schema-collection
cd schema-collection
```

6. Add a picture of a fish and a picture of a dog to a subdirectory called `subdir1`:

```shell
mkdir subdir1

cp ${EMPOROUS_CLIENT_GO_REPO}/test/fish.jpg .
cp ${EMPOROUS_CLIENT_GO_REPO}/cli/testdata/emporous-template/dog.jpeg subdir1/dog.jpg
```

7. Create a json document describing the two resources created within the workspace in a file called `schema-collection.json`

```bash
cat << EOF > schema-collection.json
{
    "fish": "fish.jpg",
    "dog": "subdir1/dog.jpg",
}
EOF
```

8. Navigate up one directory so that the Dataset Configuration file can be created:

```shell
cd ..
```

Create the `dataset-config.yaml` for the collection with the following content. Notice that the schema previously published is referenced in the `schemaAddress` property:

```bash
cat << EOF > dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  schemaAddress: "localhost:5000/exercises/myschema:latest"
  files:
    - file: "fish.jpg"
      attributes:
        animal: "fish"
        habitat: "ocean"
        size: "small"
        color: "blue"
        mammal: "false"
    - file: "subdir1/dog.jpg"
      attributes:
        animal: "dog"
        habitat: "house"
        size: "medium"
        color: "brown"
        mammal: "true"
EOF
```

9. Use the emporous client _build_ subcommand referencing the dataset config, the content directory, and the destination registry location. The attributes specified will be validated against the schema provided.

```shell
emporous build collection schema-collection --plain-http localhost:5000/exercises/schemacollection:latest --dsconfig dataset-config.yaml 
```

A validation error occurred since the _mammal_ attribute in the Dataset Configuration specified a string value instead of a boolean as defined in the schema.

In order to be able to build the schema, modify the _mammal_ attribute of the `dataset-config.yaml` file by removing the surrounding quotes as shown below in the updated Dataset Configuration:

```bash
cat << EOF > dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  schemaAddress: "localhost:5000/exercises/myschema:latest"
  files:
    - file: "fish.jpg"
      attributes:
        animal: "fish"
        habitat: "ocean"
        size: "small"
        color: "blue"
        mammal: false
    - file: "subdir1/dog.jpg"
      attributes:
        animal: "dog"
        habitat: "house"
        size: "medium"
        color: "brown"
        mammal: true
EOF
```

With a valid Dataset Configuration now in place, the collection should build successfully:

```shell
emporous build collection schema-collection --plain-http localhost:5000/exercises/schemacollection:latest --dsconfig dataset-config.yaml
```

10. Use the emporous client _push_ subcommand to publish the collection to the remote repository

```shell
emporous push --plain-http localhost:5000/exercises/schemacollection:latest
```

11. Inspect the OCI manifest of the published dataset

```shell
curl -s -H "Accept: application/vnd.oci.image.manifest.v1+json" http://localhost:5000/v2/exercises/schemacollection/manifests/latest | jq -r
```

Note that the schema id is recorded in the attribute annotation within the manifest:

```json
...
"annotations": {
  "emporous.attributes":  "emporous.attributes": "{\"myschemaid\"{\"animal\":\"dog\",\"color\":\"brown\",\"habitat\":\"house\",\"mammal\":true,\"size\":\"medium\"}}"
}
...
```

With the collection pushed, you can also revisit some of the other ways of interacting with the collection including _inspecting_ the local cache or _pulling_ either the entire collection or a subset by defining an _AttributeQuery_. Refer back to the _basic_ exercise for examples of how these steps can be achieved.

### Collection Publishing with Links

Collections can also refer to other collection; known as _Linked Collections_. It is important to note that a Linked Collection **must** have an attached schema.

Be sure that the `EMPOROUS_CLIENT_GO_REPO` environment variable is defined as described at the beginning of the exercises along with the `examples` directory.

1. Create a new directory called `linked` within the _exercises_ directory and change into this directory

```shell
mkdir -p ${EMPOROUS_CLIENT_GO_REPO}/exercises/linked
cd ${EMPOROUS_CLIENT_GO_REPO}/exercises/linked
```

2. Create a `schema-config.yaml` file with the following contents to define the schema for the collection:

```bash
cat << EOF > animalschema.json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "animal": {
      "type": "string"
    },
    "size": {
      "type": "string"
    },
    "color": {
      "type": "string"
    },
    "habitat": {
      "type": "string"
    },
    "mammal": {
      "type": "boolean"
    }
  },
  "required": [
    "animal",
    "size",
    "color",
    "habitat",
    "mammal"
  ]
}
EOF
```
```bash
cat << EOF > schema-config.yaml
kind: SchemaConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
schema:
  id: myschemaid
  schemaPath: animalschema.json
EOF
```

3. Build and push the schema to the remote registry

```bash
emporous build schema schema-config.yaml localhost:5000/exercises/linkedschema:latest
emporous push --plain-http localhost:5000/exercises/linkedschema:latest
```

4. To demonstrate Linked Collections, start creating a leaf collection by creating a workspace directory called `leaf-workspace`.

```bash
mkdir leaf-workspace
```

5. Create a simple file called `leaf.txt` containing a single word within the workspace:

```bash
echo "leaf" > leaf-workspace/leaf.txt
```

6. Create the Dataset Configuration within a file called `leaf-dataset-config.yaml` with the following content:

```bash
cat << EOF > leaf-dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  schemaAddress: localhost:5000/exercises/linkedschema:latest
  files:
    - file: "*.txt"
      attributes:
        animal: "fish"
        habitat: "ocean"
        size: "small"
        color: "blue"
        type: "leaf"
EOF
```

7. Build and push the leaf collection to the remote registry

```bash
emporous build collection leaf-workspace --plain-http localhost:5000/exercises/leaf:latest --dsconfig leaf-dataset-config.yaml
emporous push --plain-http localhost:5000/exercises/leaf:latest
```

8. Build a Root collection and link the previously built collection

Create a new directory for the root collection called `root-workspace`

```bash
mkdir root-workspace
```

9. Create a simple file called `root.txt` containing a single word within the workspace:

```bash
echo "root" > root-workspace/root.txt
```

10. Create the Dataset Configuration within a file called `root-dataset-config.yaml` with the following content:

```bash
cat << EOF > root-dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  linkedCollections:
  - localhost:5000/exercises/leaf:latest
  schemaAddress: localhost:5000/exercises/linkedschema:latest
  files:
    - file: "*.txt"
      attributes:
        animal: "cat"
        habitat: "house"
        size: "small"
        color: "orange"
        type: "root"
EOF
```

11. Build and push the root collection to the remote registry

```bash
emporous build collection root-workspace --plain-http localhost:5000/exercises/root:latest --dsconfig root-dataset-config.yaml
emporous push --plain-http localhost:5000/exercises/root:latest
```

12. Pull the collection into a directory called `linked-output`

```bash
emporous pull --plain-http localhost:5000/exercises/root:latest -o linked-output
```

13. Inspect the contents of the `linked-output` directory

```bash
ls linked-output

root.txt
```

14. Retrieve the contents of the root and leaf collection

Notice that only the contents of the root collection was retrieved in the prior step. To pull the content of both the root and any leaf collections, use the `--pull-all` flag of the emporous client `pull` subcommand into a directory called `all-linked-output`.

```bash
emporous pull --plain-http localhost:5000/exercises/root:latest --pull-all -o all-linked-output
```

15. Inspect the content of the `all-linked-output` directory

```bash
ls all-linked-output

leaf.txt root.txt
```

16. Pulling by attributes can also be specified when referencing Linked Collections. Create a file called `color-query.yaml` to content that has the attribute `color=orange`

```bash
cat << EOF > color-query.yaml
kind: AttributeQuery
apiVersion: client.emporous-framework.io/v1alpha1
attributes:
  "myschemaid":
     "color": "orange"
EOF
```

17. Pull the contents from the linked collection using the Attribute Query into a directory called `color-output`:

```bash
emporous pull --plain-http localhost:5000/exercises/root:latest --pull-all --attributes color-query.yaml -o color-output
```

18. Inspect the content of the `color-output` directory

```bash
ls color-output

root.txt
```

Notice how only the _root.txt_ file was retrieved as only this file contained the attribute `color=orange`

# Experimental

## Publish content to use with a container runtime
#### Steps
1. Create a simple Go application
```bash
cat << EOF > main.go
package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
EOF
go build -o myworkspace/helloworld main.go 
```
2. Create a collection
```bash
cat << EOF > dataset-config.yaml
kind: DataSetConfiguration
apiVersion: client.emporous-framework.io/v1alpha1
collection:
  runtime:
    Cmd:
      - "./helloworld"
  files:
    - file: "helloworld"
      fileInfo:
        permissions: 0700
      attributes:
        test: "something"
EOF
```
```bash
emporous build collection myworkspace --plain-http localhost:5000/exercises/runtime:latest --dsconfig dataset-config.yaml
emporous push --plain-http localhost:5000/exercises/runtime:latest
```

# Glossary

`collection`: a collection of linked files represented as on OCI artifact
`schema`: the properties and datatypes that can be specified within a collection

# Schemas
```bash
# core-link
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "registryHint": {
      "type": "string"
    },
    "namespaceHint": {
      "type": "string"
    },
    "transitive": {
      "type": "boolean"
    }
  },
  "required": [
    "registryHint",
    "namespaceHint",
    "transitive"
  ]
}
```
```bash
#core-descriptor
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "version": {
      "type": "string"
    },
    "type": {
      "type": "string"
    },
    "foundBy": {
      "type": "string"
    },
    "locations": {
      "type": "null"
    },
    "licenses": {
      "type": "null"
    },
    "language": {
      "type": "string"
    },
    "cpes": {
      "type": "null"
    },
    "purl": {
      "type": "string"
    }
  },
  "required": [
    "id",
    "name",
    "version",
    "type",
    "foundBy",
    "locations",
    "licenses",
    "language",
    "cpes",
    "purl"
  ]
}
```
```bash
#core-schema
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    }
  },
  "required": [
    "id"
  ]
}
```
```bash
#core-runtime


{
  "description": "OpenContainer Config Specification",
  "$schema": "https://json-schema.org/draft-04/schema#",
  "id": "https://opencontainers.org/schema/image/config",
  "type": "object",
  "properties": {
    "created": {
      "type": "string",
      "format": "date-time"
    },
    "author": {
      "type": "string"
    },
    "architecture": {
      "type": "string"
    },
    "variant": {
      "type": "string"
    },
    "os": {
      "type": "string"
    },
    "os.version": {
      "type": "string"
    },
    "os.features": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "config": {
      "type": "object",
      "properties": {
        "User": {
          "type": "string"
        },
        "ExposedPorts": {
          "$ref": "defs.json#/definitions/mapStringObject"
        },
        "Env": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "Entrypoint": {
          "oneOf": [
            {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            {
              "type": "null"
            }
          ]
        },
        "Cmd": {
          "oneOf": [
            {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            {
              "type": "null"
            }
          ]
        },
        "Volumes": {
          "oneOf": [
            {
              "$ref": "defs.json#/definitions/mapStringObject"
            },
            {
              "type": "null"
            }
          ]
        },
        "WorkingDir": {
          "type": "string"
        },
        "Labels": {
          "oneOf": [
            {
              "$ref": "defs.json#/definitions/mapStringString"
            },
            {
              "type": "null"
            }
          ]
        },
        "StopSignal": {
          "type": "string"
        },
        "ArgsEscaped": {
          "type": "boolean"
        }
      }
    },
    "rootfs": {
      "type": "object",
      "properties": {
        "diff_ids": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "type": {
          "type": "string",
          "enum": [
            "layers"
          ]
        }
      },
      "required": [
        "diff_ids",
        "type"
      ]
    },
    "history": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "created": {
            "type": "string",
            "format": "date-time"
          },
          "author": {
            "type": "string"
          },
          "created_by": {
            "type": "string"
          },
          "comment": {
            "type": "string"
          },
          "empty_layer": {
            "type": "boolean"
          }
        }
      }
    }
  },
  "required": [
    "architecture",
    "os",
    "rootfs"
  ]
}
```
```bash
#core-file
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "permissions": {
      "type": "integer"
    },
    "uid": {
      "type": "integer"
    },
    "gid": {
      "type": "integer"
    }
  },
  "required": [
    "permissions",
    "uid",
    "gid"
  ]
}
```