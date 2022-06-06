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



