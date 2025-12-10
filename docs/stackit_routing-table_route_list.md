## stackit routing-table route list

Lists all routes within a routing-table

### Synopsis

Lists all routes within a routing-table

```
stackit routing-table route list [flags]
```

### Examples

```
  List all routes within a routing-table
  $ stackit routing-table route list --routing-table-id xxx --organization-id yyy --network-area-id zzz

  List all routes within a routing-table with labels
  $ stackit routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc

  List all routes within a routing-tables with labels and limit to 10
  $ stackit routing-table list --routing-table-id xxx --organization-id yyy --network-area-id zzz --label-selector env=dev,env=rc --limit 10
```

### Options

```
  -h, --help                      Help for "stackit routing-table route list"
      --label-selector string     Filter by label
      --limit int                 Maximum number of entries to list
      --network-area-id string    Network-Area ID
      --organization-id string    Organization ID
      --routing-table-id string   Routing-Table ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit routing-table route](./stackit_routing-table_route.md)	 - Manages routes of a routing-table

