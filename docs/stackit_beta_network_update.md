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

  Update network with ID "xxx" with new name "network-1-new" and new dns servers
  $ stackit beta network update xxx --name network-1-new --dns-servers "2.2.2.2"
```

### Options

```
      --dns-servers strings   List of DNS servers/nameservers IPs
  -h, --help                  Help for "stackit beta network update"
  -n, --name string           Network name
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

