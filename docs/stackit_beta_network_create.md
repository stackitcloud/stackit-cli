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

  Create a network with name "network-1" with DNS name servers and a prefix length
  $ stackit beta network create --name network-1  --dns-name-servers "1.1.1.1,8.8.8.8,9.9.9.9" --prefix-length 25
```

### Options

```
      --dns-name-servers strings   List of DNS name servers IPs
  -h, --help                       Help for "stackit beta network create"
  -n, --name string                Network name
      --prefix-length int          The prefix length of the network
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

