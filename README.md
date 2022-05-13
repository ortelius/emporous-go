# Universal Runtime (UR) Client
Note: The UR Client is being actively developed. Please consider joining the UOR Community to participate!

## Participate
Please join us in the discussion forum and feel free to ask questions about the UOR-Framework or UR Client.

## About
The Universal Runtime Client interacts with UOR artifacts and is aware of the runtime instruction
embedded in UOR artifacts.

To learn more about Universal Runtime visit the UOR Framework website at https://uor-framework.github.io.

## Basic Usage

### Template content in a directory without pushing 
```
# The default workspace is "client-workspace" in the current working directory
client build <directory> --output my-workspace
```

### Template content in a directory and push to a registry location
`client build <directory> --push --destination localhost:5000/myartifacts:latest`



