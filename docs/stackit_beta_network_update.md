## stackit beta network update

Updates a network

### Synopsis

Updates a network.

```
stackit beta network update [flags]
```

### Examples

```
  Update network with ID "xxx" with new name "network-1-new"
  $ stackit beta network update xxx --name network-1-new

  Update network with ID "xxx" with no gateway
  $ stackit beta network update --no-ipv4-gateway

  Update IPv4 network with ID "xxx" with new name "network-1-new", new gateway and new DNS name servers
  $ stackit beta network update xxx --name network-1-new --ipv4-dns-name-servers "2.2.2.2" --ipv4-gateway "10.1.2.3"

  Update IPv6 network with ID "xxx" with new name "network-1-new", new gateway and new DNS name servers
  $ stackit beta network update xxx --name network-1-new --ipv6-dns-name-servers "2001:4860:4860::8888" --ipv6-gateway "2001:4860:4860::8888"
```

### Options

```
  -h, --help                            Help for "stackit beta network update"
      --ipv4-dns-name-servers strings   List of DNS name servers IPv4. Nameservers cannot be defined for routed networks
      --ipv4-gateway string             The IPv4 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway
      --ipv6-dns-name-servers strings   List of DNS name servers for IPv6. Nameservers cannot be defined for routed networks
      --ipv6-gateway string             The IPv6 gateway of a network. If not specified, the first IP of the network will be assigned as the gateway
  -n, --name string                     Network name
      --no-ipv4-gateway                 If set to true, the network doesn't have an IPv4 gateway
      --no-ipv6-gateway                 If set to true, the network doesn't have an IPv6 gateway
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta network](./stackit_beta_network.md)	 - Provides functionality for Network

