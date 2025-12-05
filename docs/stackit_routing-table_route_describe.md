## stackit routing-table route describe

Describes a route within a routing-table

### Synopsis

Describes a route within a routing-table

```
stackit routing-table route describe ROUTE_ID [flags]
```

### Examples

```
  Describe a route within a routing-table
  $ stackit routing-table route describe xxx --routing-table-id xxx --organization-id yyy --network-area-id zzz
```

### Options

```
  -h, --help                      Help for "stackit routing-table route describe"
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

