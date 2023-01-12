## emporous serve

Serve gRPC API to allow emporous collection management

```
emporous serve SOCKET [flags]
```

### Examples

```
  # Serve with a specified unix domain socket location
  emporous serve /var/run/test.sock
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

* [emporous](emporous.md)	 - emporous Client

