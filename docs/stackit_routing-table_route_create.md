## stackit routing-table route create

Creates a route in a routing-table

### Synopsis

Creates a route in a routing-table.

```
stackit routing-table route create [flags]
```

### Examples

```
  Create a route with CIDRv4 destination and IPv4 nexthop
  stackit routing-table route create  \ 
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv4 --destination-value <ipv4-cidr> \
--nexthop-type ipv4 --nexthop-value <ipv4-address>

  Create a route with CIDRv6 destination and IPv6 nexthop
  stackit routing-table route create \
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv6 --destination-value <ipv6-cidr> \
--nexthop-type ipv6 --nexthop-value <ipv6-address>

  Create a route with CIDRv6 destination and Nexthop Internet
  stackit routing-table route create \
--routing-table-id xxx --organization-id yyy --network-area-id zzz \
--destination-type cidrv6 --destination-value <ipv6-cidr> \
--nexthop-type internet
```

### Options

```
      --destination-type string    Destination type
      --destination-value string   Destination value
  -h, --help                       Help for "stackit routing-table route create"
      --labels stringToString      Key=value labels (default [])
      --network-area-id string     Network-Area ID
      --nexthop-type string        Next hop type
      --nexthop-value string       NextHop value
      --organization-id string     Organization ID
      --routing-table-id string    Routing-Table ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit routing-table route](./stackit_routing-table_route.md)	 - Manages routes of a routing-table

