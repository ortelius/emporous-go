## uor-client-go push

Push a UOR collection into a registry

```
uor-client-go push DST [flags]
```

### Examples

```
  # Push artifacts.
  uor-client-go push localhost:5000/myartifacts:latest
```

### Options

```
  -c, --configs stringArray   Path(s) to your registry credentials. Defaults to well-known auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.
  -h, --help                  help for push
      --insecure              Allow connections to registries SSL registry without certs
      --plain-http            Use plain http and not https when contacting registries
  -s, --sign                  keyless OIDC signing of UOR Collections with Sigstore
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [uor-client-go](uor-client-go.md)	 - UOR Client

