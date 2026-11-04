## stackit alpha vpc routing-table describe

Describes a VPC routing-table

### Synopsis

Describes a routing-table for a VPC in a region.

```
stackit alpha vpc routing-table describe ROUTING_TABLE_ID [flags]
```

### Examples

```
  Describe routing-table "xxx" for VPC "yyy" in region "eu02"
  $ stackit alpha vpc routing-table describe xxx --vpc-id yyy --region eu02

  Describe a routing-table in JSON format, using the configured project and region
  $ stackit alpha vpc routing-table describe xxx --vpc-id yyy --output-format json
```

### Options

```
  -h, --help            Help for "stackit alpha vpc routing-table describe"
      --vpc-id string   VPC ID
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

* [stackit alpha vpc routing-table](./stackit_alpha_vpc_routing-table.md)	 - Manages VPC routing-tables

