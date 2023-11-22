## stackit ske cluster create

Creates an SKE cluster

### Synopsis

Creates an SKE cluster

```
stackit ske cluster create [flags]
```

### Examples

```
$ stackit ske cluster create --project-id xxx --payload @./payload.json
```

### Options

```
  -h, --help             help for create
  -n, --name string      Cluster name
      --payload string   Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit ske cluster](./stackit_ske_cluster.md)	 - Provides functionality for SKE cluster

