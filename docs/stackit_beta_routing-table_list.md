## stackit beta routing-table list

List all routing-tables

### Synopsis

List all routing-tables

```
stackit beta routing-table list [flags]
```

### Examples

```
  List all routing-tables
  $ stackit beta routing-table list --organization-id xxx --network-area-id yyy

  List all routing-tables with labels
  $ stackit beta routing-table list --label-selector env=dev,env=rc --organization-id xxx --network-area-id yyy

  List all routing-tables with labels and set limit to 10
  $ stackit beta routing-table list --label-selector env=dev,env=rc --limit 10 --organization-id xxx --network-area-id yyy
```

### Options

```
  -h, --help                     Help for "stackit beta routing-table list"
      --label-selector string    Filter by label
      --limit int                Maximum number of entries to list
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

