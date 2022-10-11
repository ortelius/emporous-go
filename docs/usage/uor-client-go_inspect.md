## uor-client-go inspect

Print UOR collection information

```
uor-client-go inspect [flags]
```

### Examples

```
  # List all references
  uor-client-go inspect
  
  # List all descriptors for reference
  uor-client-go inspect --reference localhost:5001/test:latest
  
  # List all descriptors for reference with attribute filtering
  uor-client-go inspect --reference localhost:5001/test:latest --attributes attribute-query.yaml
```

### Options

```
  -a, --attributes string   Attribute query config path
  -h, --help                help for inspect
  -r, --reference string    A reference to list descriptors for
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [uor-client-go](uor-client-go.md)	 - UOR Client

