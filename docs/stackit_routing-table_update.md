## stackit routing-table update

Updates a routing-table

### Synopsis

Updates a routing-table.

```
stackit routing-table update ROUTING_TABLE_ID [flags]
```

### Examples

```
  Updates the label(s) of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"
  $ stackit routing-table update xxx --labels key=value,foo=bar --organization-id yyy --network-area-id zzz

  Updates the name of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"
  $ stackit routing-table update xxx --name foo --organization-id yyy --network-area-id zzz

  Updates the description of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"
  $ stackit routing-table update xxx --description foo --organization-id yyy --network-area-id zzz

  Disables the dynamic routes of a routing-table with ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"
  $ stackit routing-table update xxx --organization-id yyy --network-area-id zzz --non-dynamic-routes
```

### Options

```
      --description string       Description of the routing-table
  -h, --help                     Help for "stackit routing-table update"
      --labels stringToString    Key=value labels (default [])
      --name string              Name of the routing-table
      --network-area-id string   Network-Area ID
      --non-dynamic-routes       If true, preventing dynamic routes from propagating to the routing-table.
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

