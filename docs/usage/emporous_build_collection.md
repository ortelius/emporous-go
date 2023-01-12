## emporous build collection

Build and save an OCI artifact from files

```
emporous build collection SRC DST [flags]
```

### Examples

```
  # Build artifacts.
  emporous build collection my-directory localhost:5000/myartifacts:latest
  
  # Build artifacts with custom annotations.
  emporous build collection my-directory localhost:5000/myartifacts:latest --dsconfig dataset-config.yaml
```

### Options

```
  -c, --configs stringArray   Path(s) to your registry credentials. Defaults to well-known auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.
  -d, --dsconfig string       config path for artifact building and dataset configuration
  -h, --help                  help for collection
      --insecure              Allow connections to registries SSL registry without certs
      --no-verify             skip schema signature verification
      --plain-http            Use plain http and not https when contacting registries
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [emporous build](emporous_build.md)	 - Build and save an OCI artifact from files

