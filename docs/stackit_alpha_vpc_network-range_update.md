## stackit alpha vpc network-range update

Updates a regional network range in a VPC

### Synopsis

Updates a regional network range in a VPC.

```
stackit alpha vpc network-range update NETWORK_RANGE_ID [flags]
```

### Examples

```
  Update a network range with ID "xxx" in a VPC with ID "yyy" and with ip version "ipv4" to new nameservers "1.2.3.4,5.6.7.8"
  $ stackit alpha vpc network-range update xxx --vpc-id yyy --ip-version ipv4 --nameservers "1.2.3.4,5.6.7.8"

  Update a network range with ID "xxx" in a VPC with ID "yyy" and with ip version "ipv4" to new description "updated network range" and default prefix length 24
  $ stackit alpha vpc network-range update xxx --vpc-id yyy --ip-version ipv4 --description "updated network range" --default-prefix-length 24
```

### Options

```
      --default-prefix-length int   The default prefix length for network ranges in the VPC
      --description string          Description of the network range
  -h, --help                        Help for "stackit alpha vpc network-range update"
      --ip-version string           IP version of the network-range (one of: [ipv4])
      --labels stringToString       Labels are key-value string pairs which can be attached to a network range. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --max-prefix-length int       The maximal prefix length for network ranges in the VPC
      --min-prefix-length int       The minimal prefix length for network ranges in the VPC
      --nameservers strings         A list containing DNS Servers
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

