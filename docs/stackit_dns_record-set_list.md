## stackit dns record-set list

List all DNS record sets

### Synopsis

List all DNS record sets

```
stackit dns record-set list [flags]
```

### Examples

```
$ stackit dns record-set list --project-id xxx --zone-id xxx
```

### Options

```
  -h, --help                   help for list
      --is-active              Filter by active status, one of ["true" "false"]
      --name-like string       Filter by name
      --order-by-name string   Order by name, one of ["asc" "desc"]
      --zone-id string         Zone ID
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

