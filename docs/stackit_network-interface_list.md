## stackit network-interface list

Lists all network interfaces of a network

### Synopsis

Lists all network interfaces of a network.

```
stackit network-interface list [flags]
```

### Examples

```
  Lists all network interfaces
  $ stackit network-interface list

  Lists all network interfaces with network ID "xxx"
  $ stackit network-interface list --network-id xxx

  Lists all network interfaces with network ID "xxx" which contains the label xxx
  $ stackit network-interface list --network-id xxx --label-selector xxx

  Lists all network interfaces with network ID "xxx" in JSON format
  $ stackit network-interface list --network-id xxx --output-format json

  Lists up to 10 network interfaces with network ID "xxx"
  $ stackit network-interface list --network-id xxx --limit 10
```

### Options

```
  -h, --help                    Help for "stackit network-interface list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
      --network-id string       Network ID
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

* [stackit network-interface](./stackit_network-interface.md)	 - Provides functionality for network interfaces

