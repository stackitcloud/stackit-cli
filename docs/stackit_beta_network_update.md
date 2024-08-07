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

  Update IPv4 network with ID "xxx" with new name "network-1-new" and new DNS name servers
  $ stackit beta network update xxx --name network-1-new --ipv4-dns-name-servers "2.2.2.2"

  Update IPv6 network with ID "xxx" with new name "network-1-new" and new DNS name servers
  $ stackit beta network update xxx --name network-1-new --ipv6-dns-name-servers "2.2.2.2"
```

### Options

```
  -h, --help                            Help for "stackit beta network update"
      --ipv4-dns-name-servers strings   List of DNS name servers IPv4
      --ipv6-dns-name-servers strings   List of DNS name servers for IPv6
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

