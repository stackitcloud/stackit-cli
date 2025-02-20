## stackit server network-interface detach

Detaches a network interface from a server

### Synopsis

Detaches a network interface from a server.

```
stackit server network-interface detach [flags]
```

### Examples

```
  Detach a network interface with ID "xxx" from a server with ID "yyy"
  $ stackit server network-interface detach --network-interface-id xxx --server-id yyy

  Detach and delete all network interfaces for network with ID "xxx" and detach them from a server with ID "yyy"
  $ stackit server network-interface detach --network-id xxx --server-id yyy --delete
```

### Options

```
  -b, --delete                        If this is set all network interfaces will be deleted. (default false)
  -h, --help                          Help for "stackit server network-interface detach"
      --network-id string             Network ID
      --network-interface-id string   Network Interface ID
      --server-id string              Server ID
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

* [stackit server network-interface](./stackit_server_network-interface.md)	 - Allows attaching/detaching network interfaces to servers

