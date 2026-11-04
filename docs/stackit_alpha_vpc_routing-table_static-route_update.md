## stackit alpha vpc routing-table static-route update

Updates a static-route in a VPC routing-table

### Synopsis

Updates a static-route in a VPC routing-table.

```
stackit alpha vpc routing-table static-route update ROUTE_ID [flags]
```

### Examples

```
  Updates the label(s) of a static-route with ID "xxx" in a routing-table ID "yyy" in VPC "zzz"
  $ stackit alpha vpc routing-table static-route update xxx --routing-table-id yyy --vpc-id zzz --labels key=value,foo=bar
```

### Options

```
  -h, --help                      Help for "stackit alpha vpc routing-table static-route update"
      --labels stringToString     Labels are key-value string pairs which can be attached to a static route. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --routing-table-id string   Routing-Table ID
      --vpc-id string             VPC ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit alpha vpc routing-table static-route](./stackit_alpha_vpc_routing-table_static-route.md)	 - Manages static routes of a VPC routing-table

