## stackit routing-table delete

Deletes a routing-table

### Synopsis

Deletes a routing-table

```
stackit routing-table delete ROUTING_TABLE [flags]
```

### Examples

```
  Deletes a a routing-table
  $ stackit routing-table delete xxx --organization-id yyy --network-area-id zzz
```

### Options

```
  -h, --help                     Help for "stackit routing-table delete"
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

