## stackit alpha vpc routing-table static-route describe

Describes a static-route within a VPC routing-table

### Synopsis

Describes a static-route within a VPC routing-table.

```
stackit alpha vpc routing-table static-route describe ROUTE_ID [flags]
```

### Examples

```
  Describe a static-route within a VPC routing-table
  $ stackit alpha vpc routing-table static-route describe xxx --routing-table-id yyy --vpc-id zzz
```

### Options

```
  -h, --help                      Help for "stackit alpha vpc routing-table static-route describe"
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

