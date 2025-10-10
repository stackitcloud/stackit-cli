## stackit beta routing-table describe

Describe a routing-table

### Synopsis

Describe a routing-table

```
stackit beta routing-table describe ROUTING_TABLE_ID_ARG [flags]
```

### Examples

```
  Describe a routing-table
  $ stackit beta routing-table describe xxxx-xxxx-xxxx-xxxx --organization-id xxx --network-area-id yyy
```

### Options

```
  -h, --help                     Help for "stackit beta routing-table describe"
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

* [stackit beta routing-table](./stackit_beta_routing-table.md)	 - Manage routing-tables and its according routes

