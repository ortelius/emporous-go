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
```
# This is a optional precursor to building a collection. Collections can reference an
# already built schema or no schema at all.
uor-client-go build schema schema-config.yaml localhost:5000/myschema:latest
```
### Build workspace into an artifact

```
uor-client-go build collection my-workspace localhost:5000/myartifacts:latest
```

```
# Optionally with dsconfig
uor-client-go build my-workspace localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml
```
### Push workspace to a registry location

```
uor-client-go push my-workspace localhost:5000/myartifacts:latest
```

### Pull UOR collection to a location

```
uor-client-go pull localhost:5000/myartifacts:latest -o my-output-directory
```

### Pull subsets of a UOR collection to a location by attribute

```
uor-client-go pull localhost:5000/myartifacts:latest -o my-output-directory --attributes key=value
```

## Getting Started

### Basic Collection Publishing
1. Create a new directory.
2. Add the content to be uploaded in the directory (can be files of any content types).
3. Create a json doc where the value of each kv pair is the path to each file within the directory. Multiple json docs can be used to create deep graphs, but a graph must only have one root. Multiple json docs in a build directory is for advanced use cases. Most use cases do not need more than one json doc.

    Example json doc:
    
    ```bash
    {
        "fish": "fish.jpg",
        "text": "subdir1/file.txt",
        "fish2": "subdir1/fish2.jpg"
    }
    ```

4. Create a dataset-config.yaml outside the content directory that references the relative paths from within the content directory to each file. Add user defined key value pairs as subkeys to the `annotations`section. Each file should have as many attributes as possible. Multiple files can be referenced by using the `*` wildcard.

    Example dataset-config.yaml:
    
    ```bash
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

5. Run the UOR client build command referencing the dataset config, the content directory, and the destination registry location.
    ```
    uor-client-go build my-workspace localhost:5000/test/dataset:latest --dsconfig dataset-config.yaml 
    ```
6. Run the UOR push command to publish
    ```
    uor-client-go push localhost:5000/test/dataset:latest
    ```

7. Optionally inspect the OCI manifest of the dataset:
  `curl -H "Accept: application/vnd.oci.image.manifest.v1+json" <servername>:<port>/v2/<namespace>/<repo>/manifests/<digest or tag>`

8. Optionally inspect the cache:
      `uor-client-go inspect`

9. Optionally pull the collection back down to verify the content with `uor-client-go pull`:
      `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory`

10. Optionally pull a subset of the collection back down to verify the content with `uor-client-go pull`:

     Example attribute-query.yaml:
     ```bash
     kind: AttributeQuery
     apiVersion: client.uor-framework.io/v1alpha1
     attributes:
       fiction: true
     ```
     `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory --attributes attribute-query.yaml`

### Collection Publishing with Schema
1. Create a schema-configuration file to define attribute keys and types for corresponding collections:
   Example schema-config.yaml
   
   ```bash
    kind: SchemaConfiguration
    apiVersion: client.uor-framework.io/v1alpha1
    schema:
      attributeTypes:
        "animal": string
        "size": number
        "color": string
        "habitat": string
        "mammal": boolean
   ```
2. Build and save the schema:
    ```
    uor-client-go build schema schema-config.yaml localhost:5000/myschema:latest
    ```
3. Push the schema to the remote registry:
   ```
   uor-client-go push localhost:5000/myschema:latest
   ```
5. Create a new directory.
6. Add the content to be uploaded in the directory (can be files of any content types).
7. Create a json doc where the value of each kv pair is the path to each file within the directory. Multiple json docs can be used to create deep graphs, but a graph must only have one root. Multiple json docs in a build directory is for advanced use cases. Most use cases do not need more than one json doc.

   Example json doc:

    ```bash
    {
        "fish": "fish.jpg",
        "dog": "subdir1/dog.jpg",
    }
    ```

8. Create a dataset-config.yaml outside the content directory that references the relative paths from within the content directory to each file. Add user defined key value pairs as subkeys to the `annotations`section. Each file should have as many attributes as possible. Multiple files can be referenced by using the `*` wildcard.

   Example dataset-config.yaml:

    ```bash
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

9. Run the UOR client build command referencing the dataset config, the content directory, and the destination registry location. The attributes specified will be validated against the schema provided.
    ```
    uor-client-go build my-workspace localhost:5000/test/dataset:latest --dsconfig dataset-config.yaml 
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
