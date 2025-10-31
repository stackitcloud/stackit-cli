## stackit network-area route create

Creates a static route in a STACKIT Network Area (SNA)

### Synopsis

Creates a static route in a STACKIT Network Area (SNA).
This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.


```
stackit network-area route create [flags]
```

### Examples

```
  Create a static route with destination "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit network-area route create --organization-id yyy --network-area-id xxx --destination 1.1.1.0/24 --next-hop 1.1.1.1

  Create a static route with labels "key:value" and "foo:bar" with destination "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit network-area route create --labels key=value,foo=bar --organization-id yyy --network-area-id xxx --destination 1.1.1.0/24 --next-hop 1.1.1.1
```

### Options

```
      --destination string       Destination route. Must be a valid IPv4 or IPv6 CIDR
  -h, --help                     Help for "stackit network-area route create"
      --labels stringToString    Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
      --network-area-id string   STACKIT Network Area ID
      --next-hop-ipv4 string     Next hop IPv4 address
      --next-hop-ipv6 string     Next hop IPv6 address
      --nexthop-blackhole        Sets next hop to black hole
      --nexthop-internet         Sets next hop to internet
      --organization-id string   Organization ID
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

* [stackit network-area route](./stackit_network-area_route.md)	 - Provides functionality for static routes in STACKIT Network Areas

