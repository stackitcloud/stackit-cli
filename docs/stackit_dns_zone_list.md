## stackit dns zone list

List DNS zones

### Synopsis

List DNS zones. Successfully deleted zones are not listed by default.

```
stackit dns zone list [flags]
```

### Examples

```
  List DNS zones
  $ stackit dns zone list

  List DNS zones in JSON format
  $ stackit dns zone list --output-format json

  List up to 10 DNS zones
  $ stackit dns zone list --limit 10

  List the deleted DNS zones
  $ stackit dns zone list --deleted
```

### Options

```
      --active                 Filter for active zones
  -h, --help                   Help for "stackit dns zone list"
      --inactive               Filter for inactive zones
      --include-deleted        Includes successfully deleted zones (if unset, these are filtered out)
      --limit int              Maximum number of entries to list
      --name-like string       Filter by name
      --order-by-name string   Order by name, one of ["asc" "desc"]
      --page-size int          Number of items fetched in each API call. Does not affect the number of items in the command output (default 100)
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

