## uor-client-go build schema

Build and save a UOR schema as an OCI artifact

```
uor-client-go build schema CFG-PATH DST [flags]
```

### Examples

```
  # Build schema artifacts.
  uor-client-go build schema schema-config.yaml localhost:5000/myartifacts:latest
```

### Options

```
  -h, --help   help for schema
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [uor-client-go build](uor-client-go_build.md)	 - Build and save an OCI artifact from files

