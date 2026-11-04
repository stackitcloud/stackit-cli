## stackit alpha vpc routing-table create

Creates a VPC routing-table

### Synopsis

Creates a routing-table for a VPC in a region. VPC can only be used in projects with enabled VPC.

```
stackit alpha vpc routing-table create [flags]
```

### Examples

```
  Create a routing-table with name "rt" for VPC "xxx" in region "eu02"
  $ stackit alpha vpc routing-table create --vpc-id xxx --region eu02 --name rt

  Create a routing-table with labels and dynamic routes disabled, using the configured project and region
  $ stackit alpha vpc routing-table create --vpc-id xxx --name rt --labels env=dev --dynamic-routes=false
```

### Options

```
      --description string      Description of the routing-table
      --dynamic-routes          If set to false, prevents dynamic routes from propagating to the routing-table (default true)
  -h, --help                    Help for "stackit alpha vpc routing-table create"
      --labels stringToString   Key=value labels (default [])
      --name string             Name of the routing-table
      --system-routes           If set to false, disables automatic system routes (default true)
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

