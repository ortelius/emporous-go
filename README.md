# Universal Runtime (UR) Client

Note: The UR Client is being actively developed. Please consider joining the UOR Community to participate!

## Participate

Please join us in the discussion forum and feel free to ask questions about the UOR-Framework or UR Client.

## About

The Universal Runtime Client interacts with UOR artifacts and is aware of the runtime instruction
embedded in UOR artifacts.

To learn more about Universal Runtime visit the UOR Framework website at <https://uor-framework.github.io>.

> WARNING: The repository is under active development and the API is subject to change.

## Development

### Requirements

- `go` version 1.17+

### Build

```
make
./bin/uor-client-go -h
```

### Test

#### Unit

```
make test-unit
```

## Basic Usage

### Version

uor-client-go version

### User Workflow

1. Use `uor-client-go build schema` to build a schema to be used with a collection.
2. Use the `uor-client-go build collection` command to build the workspace as an OCI artifact in build-cache The default location is ~/.uor/cache. It can be set with the `UOR_CACHE` environment variable`.
3. Use the `uor-client-go push` command to publish to a registry as an OCI artifact.
4. Use the `uor-client-go pull` command to pull the artifact back to a local workspace.
5. Use the `uor-client-go inspect` command to inspect the build cache to list information about references.

### Build a schema into an artifact

A schema can be optionally created prior to building a collection. Collections can then reference an already built schema or no schema at all.

```shell
uor-client-go build schema schema-config.yaml localhost:5000/myschema:latest
```

### Build workspace into an artifact

Execute the following command to build a workspace into an an artifact:

```shell
uor-client-go build collection my-workspace localhost:5000/myartifacts:latest
```

An optional dataset configuration file can be provided using the `--dsconfig` flag:

```shell
uor-client-go build my-workspace localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml
```

### Push workspace to a registry location

Push a workspace to a remote registry

```shell
uor-client-go push my-workspace localhost:5000/myartifacts:latest
```

### Pull UOR collection to a location

Pull a collection from a remote registry:

```shell
uor-client-go pull localhost:5000/myartifacts:latest -o my-output-directory
```

### Pull subsets of a UOR collection to a location by attribute

Pull a portion of a collection by filtering for a set of attribute:

```shell
uor-client-go pull localhost:5000/myartifacts:latest -o my-output-directory --attributes key=value
```

## Getting Started

This guide will walk through a basic workflow of using the UOR Client

Set an environment variable called `UOR_CLIENT_GO_REPO` with the current working directory within the repository. It will be used when working through the scenarios below.

```shell
export UOR_CLIENT_GO_REPO=$(pwd)
```

### Basic Collection Publishing

The first example creates and publishes a basic collection to a remote repository

1. Create a new _workspace_ directory for the collection

```shell
mkdir basic-collection
cd basic-collection
```

2. Add the content to be uploaded in the directory. For example, two images of fishes an a text file.

Add the first photo to the current directory

```shell
cp ${UOR_CLIENT_GO_REPO}/test/fish.jpg .
```

Next, create a directory called `subdir1` to contain several additional files:

```shell
mkdir subdir1
```

Now, add a text file and another photo:

```shell
cp ${UOR_CLIENT_GO_REPO}/test/level1/fish2.jpg subdir1/
cp ${UOR_CLIENT_GO_REPO}/test/level1/file.txt subdir1/
```

3. Create a json document where the value of each kv pair is the path to each file within the directory. Multiple json documents can be used to create deep graphs, but a graph must only have one root. Multiple json docs in a build directory is for advanced use cases which will not be covered in this basic example and most just need one json document.

Create a json document called `basic.json` in the current directory with the following content:

```json
{
    "fish": "fish.jpg",
    "text": "subdir1/file.txt",
    "fish2": "subdir1/fish2.jpg"
}
```

4. A Dataset Configuration can be use to assign _attributes_ to the various resources within a collection. This file must be located outside of the content directory and refer to the relative paths within the content directory. Add user defined key value pairs as subkeys to the `annotations`section. Each file should have as many attributes as possible. Multiple files can be referenced by using the `*` wildcard.

Navigate up one directory and create a file called `dataset-config.yaml` with the following content:

```shell
cd ..
```

```yaml
kind: DataSetConfiguration
apiVersion: client.uor-framework.io/v1alpha1
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
```

5. Run the UOR client _build_ command referencing the dataset config, the content directory, and the destination registry location to the local cache. Each of the examples in this document will make use of a registry located at `localhost:5000`. Be sure to modify as appropriate for your environment including making use of the parameters related to insecure registries or communicating over HTTP.

```shell
uor-client-go build collection basic-collection localhost:5000/test/dataset:latest --dsconfig dataset-config.yaml 
```

6. Run the UOR _push_ command to publish the collection to the remote repository.

```
uor-client-go push localhost:5000/test/dataset:latest
```

7. Inspect the OCI manifest of the published dataset. The `jq` tool can be used to format the response to make it more readable. Once again, if the remote registry is exposed using HTTP or non trusted certificates, adjust the curl command below accordingly:

```shell
curl -H "Accept: application/vnd.oci.image.manifest.v1+json" https://localhost:5000/v2/test/dataset/manifests/latest
```

8. The UOR _inspect_ subcommand can be used to view the contents of the local cache. By default, the cache is located at `~/.uor/cache/`.

```shell
uor-client-go inspect`
```

9. The collection can be pulled from the remote registry to verify the content. Use the UOR _pull_ subcommand to a directory called _my-output-directory_ using the `-o` flag:

```shell
uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory
```

10. Instead of retrieving an entire collection, a subset can be retrieved by creating a `AttributeQuery` resource.

Create a file called `attribute-query.yaml` with the following content:

```yaml
kind: AttributeQuery
apiVersion: client.uor-framework.io/v1alpha1
attributes:
  fiction: true
```

Now use the UOR _pull_ subcommand along with the `--attributes` option:

```shell
uor-client-go pull localhost:5000/test/dataset:latest -o my-filtered-output-directory --attributes attribute-query.yaml`
```

Since the `fiction=true` attribute was associated with only the _file.txt_ file it was the only resource retrieved from the collection.

### Collection Publishing with Schema

A _Schema_ can be used to define the attributes associated with a collection along with linking multiple collections.

1. Create the Schema Configuration in a file called `schema-config.yaml` to define attribute keys and types for corresponding collections:

```yaml
kind: SchemaConfiguration
apiVersion: client.uor-framework.io/v1alpha1
schema:
  attributeTypes:
    "animal": string
    "size": string
    "color": string
    "habitat": string
    "mammal": boolean
```

2. Use the UOR _build_ subcommand to build and save the schema within the local cache:

```shell
uor-client-go build schema schema-config.yaml localhost:5000/myschema:latest
```

3. Push the schema to the remote registry:

```
uor-client-go push localhost:5000/myschema:latest
```

4. Create a new directory called `schema-collection` to contain a workspace to demonstrate how a schema can be used

```shell
mkdir schema-collection
cd schema-collection
```

5. Add a picture of a fish and a picture of a dog to a subdirectory called `subdir1`:

```shell
mkdir subdir1

cp ${UOR_CLIENT_GO_REPO}/test/fish.jpg .
cp ${UOR_CLIENT_GO_REPO}/cli/testdata/uor-template/dog.jpeg subdir1/dog.jpg
```

7. Create a json document describing the two resources created within the workspace in a file called `schema-collection.json`

```json
{
    "fish": "fish.jpg",
    "dog": "subdir1/dog.jpg",
}
```

8. Navigate up one directory so that te Dataset Configuration file can be created:

```shell
cd ..
```

Create the `dataset-config.yaml` for the collection with the following content. Notice that the schema previously published is referenced in the `schemaAddress` property:

```yaml
kind: DataSetConfiguration
apiVersion: client.uor-framework.io/v1alpha1
collection:
  schemaAddress: "localhost:5000/myschema:latest"
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
    - file: "*.jpg"
      attributes:
        custom: "customval"  
```

NOTE: This step is overwriting the content of the Dataset Configuration file created in the prior section. You are free to place the contents in a file with a different name. When doing so, be sure to specify the appropriate name when executing the _build_ subcommand of the UOR client.

9. Use the UOR client _build_ subcommand referencing the dataset config, the content directory, and the destination registry location. The attributes specified will be validated against the schema provided.

```shell
uor-client-go build collection schema-collection localhost:5000/test/dataset:latest --dsconfig dataset-config.yaml 
```

A validation error occurred since the _mammal_ attribute in the Dataset Configuration specified a string value instead of a boolean as defined in the schema.

In order to be able to build the schema, modify the _mammal_ attribute of the `dataset-config.yaml` file by removing the surrounding quotes as shown below in the updated Dataset Configuration:

```yaml
kind: DataSetConfiguration
apiVersion: client.uor-framework.io/v1alpha1
collection:
  schemaAddress: "localhost:5000/myschema:latest"
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
    - file: "*.jpg"
      attributes:
        custom: "customval"  
```

10. Run the UOR push command to publish

     ```
     uor-client-go push localhost:5000/test/dataset:latest
     ```

11. Optionally inspect the OCI manifest of the dataset:
    `curl -H "Accept: application/vnd.oci.image.manifest.v1+json" <servername>:<port>/v2/<namespace>/<repo>/manifests/<digest or tag>`

12. Optionally inspect the cache:
    `uor-client-go inspect`

13. Optionally pull the collection back down to verify the content with `uor-client-go pull`:
    `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory`

14. Optionally pull a subset of the collection back down to verify the content with `uor-client-go pull`:

    Example attribute-query.yaml:

     ```bash
     kind: AttributeQuery
     apiVersion: client.uor-framework.io/v1alpha1
     attributes:
       mammal: true
     ```

    `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory --attributes attribute-query.yaml`

### Collection Publishing with Links

> IMPORTANT: Linked collection must have an attached schema

1. Build the schema

   ```bash
   vi schema-config.yaml
   
   kind: SchemaConfiguration
   apiVersion: client.uor-framework.io/v1alpha1
   schema:
    attributeTypes:
      "animal": string
      "size": number
      "color": string
      "habitat": string
      "type": string
   ```

    ```bash
    uor-client-go build schema schema-config.yaml localhost:5000/myschema:latest
    uor-client-go push localhost:5000/myschema:latest
    ```

2. Build a leaf collection

    ```bash
    mkdir leaf-workspace
    echo "leaf" > leaf-workspace/leaf.txt
    ```

    ```bash
    vi leaf-dataset-config.yaml
    
    kind: DataSetConfiguration
    apiVersion: client.uor-framework.io/v1alpha1
    collection:
      schemaAddress: localhost:5000/myschema:latest
      files:
        - file: "*.txt"
          attributes:
            animal: "fish"
            habitat: "ocean"
            size: "small"
            color: "blue"
            type: "leaf"
    ```

    ```
    uor-client-go build leaf-workspace localhost:5000/leaf:latest --dsconfig leaf-dataset-config.yaml
    uor-client-go push localhost:5000/leaf:latest
    ```

3. Build a collection and link the previously built collection

    ```bash
    mkdir root-workspace
    echo "root" > root-workspace/root.txt
    ```

    ```bash
    vi root-dataset-config.yaml
    
    kind: DataSetConfiguration
    apiVersion: client.uor-framework.io/v1alpha1
    collection:
      linkedCollections:
      - localhost:5000/leaf:latest
      schemaAddress: localhost:5000/myschema:latest
      files:
        - file: "*.txt"
          attributes:
            animal: "cat"
            habitat: "house"
            size: "small"
            color: "orange"
            type: "root"
    ```

    ```bash
    uor-client-go build root-workspace localhost:5000/root:latest --dsconfig root-dataset-config.yaml
    uor-client-go push localhost:5000/root:latest
    ```

4. Pull the collection with the `--pull-all` flag

    ```bash
    uor-client-go pull localhost:5000/root:latest 
    ls
    root.txt
    uor-client-go pull localhost:5000/root:latest --pull-all
    ls
    leaf.txt root.txt
    ```

5. Pull all with attributes

   ```bash
    vi color-query.yaml
   
    kind: AttributeQuery
    apiVersion: client.uor-framework.io/v1alpha1
    attributes:
      "color": "orange"
   ```

    ```bash
    uor-client-go pull localhost:5000/root:latest --pull-all --attributes color-query.yaml
    ls
    root.txt
    ```

    ```bash
   vi size-query.yaml
   
   kind: AttributeQuery
   apiVersion: client.uor-framework.io/v1alpha1
   attributes:
     "size": "small"
   ```

    ```bash
    uor-client-go pull localhost:5000/root:latest --pull-all --attributes size-query.yaml
    ls
    leaf.txt root.txt
    ```

# Glossary

`collection`: a collection of linked files represented as on OCI artifact
