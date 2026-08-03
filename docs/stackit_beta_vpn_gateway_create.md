## stackit beta vpn gateway create

Creates a vpn gateway

### Synopsis

Creates a vpn gateway.

```
stackit beta vpn gateway create [flags]
```

### Examples

```
  Create a vpn gateway with name "xxx", plan "p500", policy based routing and both tunnels in availability-zone eu01-1
  $ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1

  Create a vpn gateway with the labels foo=bar and x=y
  $ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1 --label foo=bar,x=y

  Create a vpn gateway with bgp enabled, yyy as local asn and [aaa, bbb] as override advertised routes
  $ stackit beta vpn gateway create --name xxx --plan-id p500 --routing-type POLICY_BASED --availability-zone-tunnel-1 eu01-1 --availability-zone-tunnel-2 eu01-1 --bgp-local-asn yyy --bgp-override-advertised-routes aaa,bbb
```

### Options

```
      --availability-zone-tunnel-1 string            Availability Zone of Tunnel 1
      --availability-zone-tunnel-2 string            Availability Zone of Tunnel 2
      --bgp-local-asn int                            ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996)
      --bgp-override-advertised-routes stringArray   A list of IPv4 Prefixes to advertise via BGP
  -h, --help                                         Help for "stackit beta vpn gateway create"
      --labels stringToString                        Labels in key=value format, separated by commas (default [])
      --name string                                  Gateway name
      --plan-id string                               Plan ID
      --routing-type string                          Routing Type of the VPN (one of: [POLICY_BASED, ROUTE_BASED, BGP_ROUTE_BASED])
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

* [stackit beta vpn gateway](./stackit_beta_vpn_gateway.md)	 - Provides functionality for VPN gateway

