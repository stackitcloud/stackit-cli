## stackit alpha vpc routing-table static-route list

Lists all static-routes within a VPC routing-table

### Synopsis

Lists all static-routes within a VPC routing-table.

```
stackit alpha vpc routing-table static-route list [flags]
```

### Examples

```
  List all static-routes within a VPC routing-table
  $ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy

  List all static-routes within a VPC routing-table with labels foo=bar,bar=baz
  $ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy --label-selector foo=bar,bar=baz

  List up to 10 static-routes within a VPC routing-table
  $ stackit alpha vpc routing-table static-route list --routing-table-id xxx --vpc-id yyy --limit 10
```

### Options

```
  -h, --help                      Help for "stackit alpha vpc routing-table static-route list"
      --label-selector string     Filter by label
      --limit int                 Maximum number of entries to list
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

