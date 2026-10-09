## stackit alpha vpc network-range list

Lists all network ranges in a VPC

### Synopsis

Lists all network ranges in a VPC.

```
stackit alpha vpc network-range list [flags]
```

### Examples

```
  Lists all network ranges in a VPC with ID "xxx"
  $ stackit alpha vpc network-range list --vpc-id xxx

  Lists all network ranges in a VPC with ID "xxx" in JSON format
  $ stackit alpha vpc network-range list --vpc-id xxx --output-format json

  Lists up to 10 network ranges in a VPC with ID "xxx"
  $ stackit alpha vpc network-range list --vpc-id xxx --limit 10
```

### Options

```
  -h, --help            Help for "stackit alpha vpc network-range list"
      --limit int       Maximum number of entries to list
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

* [stackit alpha vpc network-range](./stackit_alpha_vpc_network-range.md)	 - Provides functionality for network ranges in VPC

