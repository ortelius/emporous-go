## uor-client-go serve

Serve gRPC API to allow UOR collection management

```
uor-client-go serve SOCKET [flags]
```

### Examples

```
  # Serve with a specified unix domain socket location
  uor-client-go serve /var/run/test.sock
```

### Options

```
  -h, --help         help for serve
      --insecure     Allow connections to registries SSL registry without certs
      --plain-http   Use plain http and not https when contacting registries
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [uor-client-go](uor-client-go.md)	 - UOR Client

