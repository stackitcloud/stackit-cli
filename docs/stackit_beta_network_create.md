## stackit beta network create

Creates a network

### Synopsis

Creates a network.

```
stackit beta network create [flags]
```

### Examples

```
  Create a network with name "network-1"
  $ stackit beta network create --name network-1

  Create an IPv4 network with name "network-1" with DNS name servers and a prefix length
  $ stackit beta network create --name network-1  --ipv4-dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --ipv4-prefix-length 25

  Create an IPv6 network with name "network-1" with DNS name servers and a prefix length
  $ stackit beta network create --name network-1  --ipv6-dns-name-servers "2001:4860:4860::8888,2001:4860:4860::8844" --ipv6-prefix-length 56
```

### Options

```
  -h, --help                            Help for "stackit beta network create"
      --ipv4-dns-name-servers strings   List of DNS name servers for IPv4
      --ipv4-prefix-length int          The prefix length of the IPv4 network
      --ipv6-dns-name-servers strings   List of DNS name servers for IPv6
      --ipv6-prefix-length int          The prefix length of the IPv6 network
  -n, --name string                     Network name
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

