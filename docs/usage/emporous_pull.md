## emporous pull

Pull a Emporous collection based on content or attribute address

```
emporous pull SRC [flags]
```

### Examples

```
  # Pull collection reference.
  emporous pull localhost:5001/test:latest
  
  # Pull collection reference and all linked references.
  emporous pull localhost:5001/test:latest --pull-all
  
  # Pull all content from reference that satisfies the attribute query.
  emporous pull localhost:5001/test:latest --attributes attribute-query.yaml
```

### Options

```
      --attributes string     Attribute query config path
  -c, --configs stringArray   Path(s) to your registry credentials. Defaults to well-known auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.
  -h, --help                  help for pull
      --insecure              Allow connections to registries SSL registry without certs
      --no-verify             Skip collection signature verification
  -o, --output string         Output location for artifacts
      --plain-http            Use plain http and not https when contacting registries
      --pull-all              Pull all linked collections
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [emporous](emporous.md)	 - Emporous Client

