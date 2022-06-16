# Universal Runtime (UR) Client
Note: The UR Client is being actively developed. Please consider joining the UOR Community to participate!

## Participate
Please join us in the discussion forum and feel free to ask questions about the UOR-Framework or UR Client.

## About
The Universal Runtime Client interacts with UOR artifacts and is aware of the runtime instruction
embedded in UOR artifacts.

To learn more about Universal Runtime visit the UOR Framework website at https://uor-framework.github.io.

## Development

### Requirements

- `go` version 1.17+

### Build

```
make
./bin/client -h
```
### Test

#### Unit:
```
make test-unit
```

## Basic Usage

### User Workflow

1. Create a directory with artifacts to publish to a registry as an OCI artifact. If the files reference each other, the client will replace the in-content linked files with the content address.
   > WARNING: Currently, only JSON is supported for link replacement.
2. Use the `client build` command to create the output workspace with the rendered content. If the files in the workspace do not contain links to each other, skip this step.
3. Use the `client push` command to publish the workspace to a registry as an OCI artifact.
### Template content in a directory without pushing 
```
# The default workspace is "client-workspace" in the current working directory
client build my-directory --output my-workspace
```

### Push workspace to a registry location
```
client push my-workspace localhost:5000/myartifacts:latest
```

### Pull artifact to a location
```
client pull localhost:5000/myartifacts:latest my-output-directory
```

## Getting Started

1. Create a new directory. 
2. Add the content to be uploaded in the directory (can be files of any content types).
3. Create a json doc where the value of each kv pair is the path to each file within the directory. Multiple json docs can be used to create deep graphs, but a graph must only have one root. Multiple json docs in a build directory is for advanced use cases. Most use cases do not need more than one json doc. 

Example json doc:
```
{
    "fish": "fish.jpg",
    "text": "subdir1/file.txt",
    "fish2": "subdir1/fish2.jpg"
}
```
4. Create a dataset-config.yaml outside of the content directory that references the relative paths from within the content directory to each file. Add user defined key value pairs as subkeys to the `annotations`section. Each file should have as many attributes as possible. Multiple files can be referenced by using the `*` wildcard. 

Example dataset-config.yaml:
```
kind: DataSetConfiguration
apiVersion: client.uor-framework.io/v1alpha1
files:
  - file: fish.jpg
    attributes:
      animal: fish
      habitat: ocean
      size: small
      color: blue
  - file: subdir1/file.txt
    attributes:
      fiction: true  
      genre: science fiction
  - file: *.jpg
    attributes:
      custom: customval

```
5. Run the UOR client build command referencing the dataset config, the content directory, and optionally push to a registry location.

```
client build --dsconfig dataset-config.yaml content-dir --output my-workspace
client push my-workspace localhost:5000/test/dataset:latest
```

6. Optionally inspect the OCI manifest of the dataset:

curl -H "Accept: application/vnd.oci.image.manifest.v1+json" <servername>:<port>/v2/<namespace>/<repo>/manifests/<digest or tag>


