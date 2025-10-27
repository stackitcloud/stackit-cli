## stackit routing-table create

Creates a routing-table

### Synopsis

Creates a routing-table.

```
stackit routing-table create [flags]
```

### Examples

```
  Create a routing-table with name `rt`
  stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt"

  Create a routing-table with name `rt` and description `some description`
  stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --description "some description"

  Create a routing-table with name `rt` with system_routes disabled
  stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --non-system-routes

  Create a routing-table with name `rt` with dynamic_routes disabled
  stackit routing-table create --organization-id xxx --network-area-id yyy --name "rt" --non-dynamic-routes
```

### Options

```
      --description string       Description of the routing-table
  -h, --help                     Help for "stackit routing-table create"
      --labels stringToString    Key=value labels (default [])
      --name string              Name of the routing-table
      --network-area-id string   Network-Area ID
      --non-dynamic-routes       If true, preventing dynamic routes from propagating to the routing-table.
      --non-system-routes        If true, automatically disables routes for project-to-project communication.
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

