## stackit routing-table describe

Describes a routing-table

### Synopsis

Describes a routing-table

```
stackit routing-table describe ROUTING_TABLE_ID_ARG [flags]
```

### Examples

```
  Describe a routing-table
  $ stackit routing-table describe xxx --organization-id xxx --network-area-id yyy
```

### Options

```
  -h, --help                     Help for "stackit routing-table describe"
      --network-area-id string   Network-Area ID
      --organization-id string   Organization ID
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

* [stackit routing-table](./stackit_routing-table.md)	 - Manage routing-tables and its according routes

