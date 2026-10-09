## stackit alpha vpc network-range describe

Shows details of a network range in a VPC

### Synopsis

Shows details of a network range in a VPC.

```
stackit alpha vpc network-range describe NETWORK_RANGE_ID [flags]
```

### Examples

```
  Show details of a network range with ID "xxx" in a VPC with ID "yyy"
  $ stackit alpha vpc network-range describe xxx --vpc-id yyy
```

### Options

```
  -h, --help            Help for "stackit alpha vpc network-range describe"
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

