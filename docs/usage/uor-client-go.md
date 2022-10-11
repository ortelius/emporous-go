## uor-client-go

UOR Client

### Synopsis

The UOR client helps you build, publish, and retrieve UOR collections as an OCI artifact.

 The workflow to publish a collection is to gather files for a collection in a directory workspace and use the build sub-command. During the build process, the tag for the remote destination is specified.

 This build action will store the collection in a build cache. This location can be specified with the UOR_CACHE environment variable. The default location is ~/.uor/cache.

 After the collection has been stored, it can be retrieved and pushed to the registry with the push sub-command.

 Collections can be retrieved from the cache or the remote location (if not stored) with the pull sub-command. The pull sub-command also allows for filtering of the collection with an attribute query configuration file.

```
uor-client-go [flags]
```

### Options

```
  -h, --help              help for uor-client-go
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [uor-client-go build](uor-client-go_build.md)	 - Build and save an OCI artifact from files
* [uor-client-go inspect](uor-client-go_inspect.md)	 - Print UOR collection information
* [uor-client-go pull](uor-client-go_pull.md)	 - Pull a UOR collection based on content or attribute address
* [uor-client-go push](uor-client-go_push.md)	 - Push a UOR collection into a registry
* [uor-client-go serve](uor-client-go_serve.md)	 - Serve gRPC API to allow UOR collection management
* [uor-client-go version](uor-client-go_version.md)	 - Print the version

