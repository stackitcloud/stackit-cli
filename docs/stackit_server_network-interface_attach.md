## stackit server network-interface attach

Attaches a network interface to a server

### Synopsis

Attaches a network interface to a server.

```
stackit server network-interface attach [flags]
```

### Examples

```
  Attach a network interface with ID "xxx" to a server with ID "yyy"
  $ stackit server network-interface attach --network-interface-id xxx --server-id yyy

  Create a network interface for network with ID "xxx" and attach it to a server with ID "yyy"
  $ stackit server network-interface attach --network-id xxx --server-id yyy --create
```

### Options

```
  -b, --create                        If this is set a network interface will be created. (default false)
  -h, --help                          Help for "stackit server network-interface attach"
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

