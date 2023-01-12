## emporous

emporous Client

### Synopsis

The emporous client helps you build, publish, and retrieve emporous collections as an OCI artifact.

 The workflow to publish a collection is to gather files for a collection in a directory workspace and use the build sub-command. During the build process, the tag for the remote destination is specified.

 This build action will store the collection in a build cache. This location can be specified with the emporous_CACHE environment variable. The default location is ~/.emporous/cache.

 After the collection has been stored, it can be retrieved and pushed to the registry with the push sub-command.

 Collections can be retrieved from the cache or the remote location (if not stored) with the pull sub-command. The pull sub-command also allows for filtering of the collection with an attribute query configuration file.

```
emporous [flags]
```

### Options

```
  -h, --help              help for emporous
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [emporous build](emporous_build.md)	 - Build and save an OCI artifact from files
* [emporous inspect](emporous_inspect.md)	 - Print Emporous collection information
* [emporous pull](emporous_pull.md)	 - Pull a Emporous collection based on content or attribute address
* [emporous push](emporous_push.md)	 - Push a emporous collection into a registry
* [emporous serve](emporous_serve.md)	 - Serve gRPC API to allow emporous collection management
* [emporous version](emporous_version.md)	 - Print the version

