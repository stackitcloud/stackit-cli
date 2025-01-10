## stackit beta server volume list

Lists all server volumes

### Synopsis

Lists all server volumes.

```
stackit beta server volume list [flags]
```

### Examples

```
  List all volumes for a server with ID "xxx"
  $ stackit beta server volume list --server-id xxx

  List all volumes for a server with ID "xxx" in JSON format
  $ stackit beta server volumes list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server volume list"
  -s, --server-id string   Server ID
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

* [stackit beta server volume](./stackit_beta_server_volume.md)	 - Provides functionality for server volumes

