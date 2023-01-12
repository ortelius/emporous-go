## emporous build schema

Build and save a emporous schema as an OCI artifact

```
emporous build schema CFG-PATH DST [flags]
```

### Examples

```
  # Build schema artifacts.
  emporous build schema schema-config.yaml localhost:5000/myartifacts:latest
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

* [emporous build](emporous_build.md)	 - Build and save an OCI artifact from files

