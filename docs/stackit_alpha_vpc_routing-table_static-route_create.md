## stackit alpha vpc routing-table static-route create

Creates a static-route in a routing-table

### Synopsis

Creates a static-route in a routing-table.

```
stackit alpha vpc routing-table static-route create [flags]
```

### Examples

```
  Create a static-route with CIDRv4 destination and IPv4 nexthop
  $ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv4 --destination-value <ipv4-cidr> --nexthop-type ipv4 --nexthop-value <ipv4-address>

  Create a route with CIDRv6 destination and IPv6 nexthop
  $ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type ipv6 --nexthop-value <ipv6-address>

  Create a route with CIDRv6 destination and nexthop Internet
  $ stackit alpha vpc routing-table static-route create --routing-table-id xxx --vpc-id yyy --destination-type cidrv6 --destination-value <ipv6-cidr> --nexthop-type internet
```

### Options

```
      --destination-type string    Destination type (one of: [cidrv4, cidrv6])
      --destination-value string   Destination value
  -h, --help                       Help for "stackit alpha vpc routing-table static-route create"
      --labels stringToString      Labels are key-value string pairs which can be attached to a static route. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --nexthop-type string        Next hop type (one of: [ipv4, ipv6, internet, blackhole])
      --nexthop-value string       NextHop value
      --routing-table-id string    Routing-Table ID
      --vpc-id string              VPC ID
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

