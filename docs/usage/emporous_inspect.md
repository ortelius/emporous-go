## emporous inspect

Print Emporous collection information

```
emporous inspect [flags]
```

### Examples

```
  # List all references
  emporous inspect
  
  # List all descriptors for reference
  emporous inspect --reference localhost:5001/test:latest
  
  # List all descriptors for reference with attribute filtering
  emporous inspect --reference localhost:5001/test:latest --attributes attribute-query.yaml
```

### Options

```
  -a, --attributes string   Attribute query config path
  -h, --help                help for inspect
  -p, --print-attributes    print descriptor attributes
  -r, --reference string    A reference to list descriptors for
```

### Options inherited from parent commands

```
  -l, --loglevel string   Log level (debug, info, warn, error, fatal) (default "info")
```

### SEE ALSO

* [emporous](emporous.md)	 - emporous Client

