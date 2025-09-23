## stackit beta routing-table route delete

Deletes a route within a routing-table

### Synopsis

Deletes a route within a routing-table

```
stackit beta routing-table route delete routing-table-id [flags]
```

### Examples

```
  Deletes a route within a routing-table
  $ stackit beta routing-table route delete xxxx-xxxx-xxxx-xxxx --routing-table-id xxx --organization-id yyy --network-area-id zzz
```

### Options

```
  -h, --help                      Help for "stackit beta routing-table route delete"
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

* [stackit beta routing-table route](./stackit_beta_routing-table_route.md)	 - Manage routes of a routing-table

