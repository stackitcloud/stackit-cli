## stackit beta server network-interface list

Lists all attached network interfaces of a server

### Synopsis

Lists all attached network interfaces of a server.

```
stackit beta server network-interface list [flags]
```

### Examples

```
  Lists all attached network interfaces of server with ID "xxx"
  $ stackit beta server network-interface list --server-id xxx

  Lists all attached network interfaces of server with ID "xxx" in JSON format
  $ stackit beta server network-interface list --server-id xxx --output-format json

  Lists up to 10 attached network interfaces of server with ID "xxx"
  $ stackit beta server network-interface list --server-id xxx --limit 10
```

### Options

```
  -h, --help               Help for "stackit beta server network-interface list"
      --limit int          Maximum number of entries to list
      --server-id string   Server ID
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

* [stackit beta server network-interface](./stackit_beta_server_network-interface.md)	 - Allows attaching/detaching network interfaces to servers

