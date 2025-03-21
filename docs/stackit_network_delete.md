## stackit network delete

Deletes a network

### Synopsis

Deletes a network.
If the network is still in use, the deletion will fail


```
stackit network delete NETWORK_ID [flags]
```

### Examples

```
  Delete network with ID "xxx"
  $ stackit network delete xxx
```

### Options

```
  -h, --help   Help for "stackit network delete"
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

