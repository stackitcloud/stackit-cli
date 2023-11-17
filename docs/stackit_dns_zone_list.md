## stackit dns zone list

List all DNS zones

### Synopsis

List all DNS zones

```
stackit dns zone list [flags]
```

### Examples

```
$ stackit dns zone list --project-id xxx
```

### Options

```
  -h, --help                   help for list
      --is-active              Filter by active status, one of ["true" "false"]
      --name-like string       Filter by name
      --order-by-name string   Order by name, one of ["asc" "desc"]
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zone

