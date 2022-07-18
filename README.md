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

1. Create a directory with artifacts to publish to a registry as an OCI artifact. If the files reference each other, the client will replace the in-content linked files with the content address.
> WARNING: Currently, only JSON is supported for link replacement.
2. Use the `uor-client-go build` command to build the workspace as an OCI artifact in build cache `default is (homedir/.uor/cache). Can be set with UOR_CACHE environment variables`.
3. Use the `uor-client-go push` command to publish to a registry as an OCI artifact.
4. Use the `uor-client-go pull` command to pull the artifact back to a local workspace.

### Build workspace into an artifact

```
client build my-workspace localhost:5000/myartifacts:latest
```

```
# Optionally with dsconfig
client build my-workspace localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml
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

7. Optionally pull the collection back down to verify the content with `uor-client-go pull`:
  `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory`

8. Optionally pull a subset of the collection back down to verify the content with `uor-client-go pull`:
  `uor-client-go pull localhost:5000/test/dataset:latest -o my-output-directory --attributes "fiction=true"`

# Glossary

`collection`: a collection of linked files represented as on OCI artifact
