## stackit network create

Creates a network

### Synopsis

Creates a network.

```
stackit network create [flags]
```

### Examples

```
  Create a network with name "network-1"
  $ stackit network create --name network-1

  Create a non-routed network with name "network-1"
  $ stackit network create --name network-1 --non-routed

  Create a network with name "network-1" and no gateway
  $ stackit network create --name network-1 --no-ipv4-gateway

  Create a network with name "network-1" and labels "key=value,key1=value1"
  $ stackit beta network create --name network-1 --labels key=value,key1=value1

  Create an IPv4 network with name "network-1" with DNS name servers, a prefix and a gateway
  $ stackit network create --name network-1  --ipv4-dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --ipv4-prefix "10.1.2.0/24" --ipv4-gateway "10.1.2.3"

  Create an IPv6 network with name "network-1" with DNS name servers, a prefix and a gateway
  $ stackit network create --name network-1  --ipv6-dns-name-servers "2001:4860:4860::8888,2001:4860:4860::8844" --ipv6-prefix "2001:4860:4860::8888" --ipv6-gateway "2001:4860:4860::8888"
```

### Options

```
  -h, --help                            Help for "stackit network create"
      --ipv4-dns-name-servers strings   List of DNS name servers for IPv4. Nameservers cannot be defined for routed networks
      --ipv4-gateway string             The IPv4 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway
      --ipv4-prefix string              The IPv4 prefix of the network (CIDR)
      --ipv4-prefix-length int          The prefix length of the IPv4 network
      --ipv6-dns-name-servers strings   List of DNS name servers for IPv6. Nameservers cannot be defined for routed networks
      --ipv6-gateway string             The IPv6 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway
      --ipv6-prefix string              The IPv6 prefix of the network (CIDR)
      --ipv6-prefix-length int          The prefix length of the IPv6 network
      --labels stringToString           Labels are key-value string pairs which can be attached to a network. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string                     Network name
      --no-ipv4-gateway                 If set to true, the network doesn't have an IPv4 gateway
      --no-ipv6-gateway                 If set to true, the network doesn't have an IPv6 gateway
      --non-routed                      If set to true, the network is not routed and therefore not accessible from other networks
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

* [stackit network](./stackit_network.md)	 - Provides functionality for networks

