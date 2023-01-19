## emporous create

Create artifacts from existing OCI artifacts

```
emporous create [flags]
```

### Options

```
  -c, --configs stringArray   Path(s) to your registry credentials. Defaults to well-known auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.
  -h, --help                  help for create
      --insecure              Allow connections to registries SSL registry without certs
      --plain-http            Use plain http and not https when contacting registries
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [emporous](emporous.md)	 - Emporous Client
* [emporous create inventory](emporous_create_inventory.md)	 - Create software inventories from UOR artifacts

