## stackit alpha vpc region list

Lists all configured regions for a VPC

### Synopsis

Lists all configured regions for a VPC.

```
stackit alpha vpc region list [flags]
```

### Examples

```
  List all configured regions for a VPC with ID "xxx"
  $ stackit alpha vpc region list --vpc-id xxx

  List all configured regions for a VPC in JSON format
  $ stackit alpha vpc region list --vpc-id xxx --output-format json
```

### Options

```
  -h, --help            Help for "stackit alpha vpc region list"
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

* [stackit alpha vpc region](./stackit_alpha_vpc_region.md)	 - Manages regional configurations of a VPC

