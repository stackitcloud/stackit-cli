## stackit routing-table route update

Updates a route in a routing-table

### Synopsis

Updates a route in a routing-table.

```
stackit routing-table route update ROUTE_ID_ARG [flags]
```

### Examples

```
  Updates the label(s) of a route with ID "xxx" in a routing-table ID "xxx" in organization with ID "yyy" and network-area with ID "zzz"
  $ stackit routing-table route update xxx --labels key=value,foo=bar --routing-table-id xxx --organization-id yyy --network-area-id zzz
```

### Options

```
  -h, --help                      Help for "stackit routing-table route update"
      --labels stringToString     Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
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

