## stackit dns record-set list

Lists DNS record sets

### Synopsis

Lists DNS record sets. Successfully deleted record sets are not listed by default.

```
stackit dns record-set list [flags]
```

### Examples

```
  List DNS record-sets for zone with ID "xxx"
  $ stackit dns record-set list --zone-id xxx

  List DNS record-sets for zone with ID "xxx" in JSON format
  $ stackit dns record-set list --zone-id xxx --output-format json

  List active DNS record-sets for zone with ID "xxx"
  $ stackit dns record-set list --zone-id xxx --is-active true

  List up to 10 DNS record-sets for zone with ID "xxx"
  $ stackit dns record-set list --zone-id xxx --limit 10

  List the deleted DNS record-sets for zone with ID "xxx"
  $ stackit dns record-set list --zone-id xxx --deleted
```

### Options

```
      --active                 Filter for active record sets
      --deleted                Filter for deleted record sets
  -h, --help                   Help for "stackit dns record-set list"
      --inactive               Filter for inactive record sets. Deleted record sets are always inactive and will be included when this flag is set
      --limit int              Maximum number of entries to list
      --name-like string       Filter by name
      --order-by-name string   Order by name, one of ["asc" "desc"]
      --page-size int          Number of items fetched in each API call. Does not affect the number of items in the command output (default 100)
      --zone-id string         Zone ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

