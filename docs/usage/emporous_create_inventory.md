## emporous create inventory

Create software inventories from UOR artifacts

```
emporous create inventory SRC [flags]
```

### Examples

```
  # Build inventory from artifacts.
  emporous create inventory localhost:5000/myartifacts:latest
```

### Options

```
  -f, --format string   software inventory format. Options are cyclonedxjson or spdx22json. Default is spdx22json
  -h, --help            help for inventory
```

### Options inherited from parent commands

```
  -c, --configs stringArray   Path(s) to your registry credentials. Defaults to well-known auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.
      --insecure              Allow connections to registries SSL registry without certs
  -l, --loglevel string       Log level (debug, info, warn, error, fatal) (default "info")
      --plain-http            Use plain http and not https when contacting registries
```

### SEE ALSO

* [emporous create](emporous_create.md)	 - Create artifacts from existing OCI artifacts

