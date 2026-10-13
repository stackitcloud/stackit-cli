## stackit alpha vpc network-range create

Creates a regional network range in a VPC

### Synopsis

Creates a regional network range in a VPC.
The vpc needs to be enabled in the specified region before creating a network range.

```
stackit alpha vpc network-range create [flags]
```

### Examples

```
  Create a network range in a VPC with ID "xxx" with prefix "1.1.1.0/24" and with ip version 4
  $ stackit alpha vpc network-range create --vpc-id xxx --prefix "1.1.1.0/24" --ip-version ipv4
```

### Options

```
      --default-prefix-length int   The default prefix length for network ranges in the VPC
      --description string          Description of the network range
  -h, --help                        Help for "stackit alpha vpc network-range create"
      --ip-version string           IP version of the network-range (one of: [ipv4])
      --labels stringToString       Labels are key-value string pairs which can be attached to a network range. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --max-prefix-length int       The maximal prefix length for network ranges in the VPC
      --min-prefix-length int       The minimal prefix length for network ranges in the VPC
      --nameservers strings         A list containing DNS Servers
      --prefix string               Network range to create in CIDR notation
      --vpc-id string               VPC ID
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

