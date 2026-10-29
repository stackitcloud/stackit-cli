## stackit alpha vpc routing-table list

Lists VPC routing-tables

### Synopsis

Lists routing-tables for a VPC in a region.

```
stackit alpha vpc routing-table list [flags]
```

### Examples

```
  List routing-tables for VPC "xxx" in region "eu02"
  $ stackit alpha vpc routing-table list --vpc-id xxx --region eu02

  List up to ten routing-tables with labels, using the configured project and region
  $ stackit alpha vpc routing-table list --vpc-id xxx --label-selector env=dev --limit 10
```

### Options

```
      --filter string           Filter resources by fields. A subset of expr-lang is supported. See https://expr-lang.org/docs/language-definition for usage details
  -h, --help                    Help for "stackit alpha vpc routing-table list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
      --vpc-id string           VPC ID
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

