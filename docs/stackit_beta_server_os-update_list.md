## stackit beta server os-update list

Lists all server os-updates

### Synopsis

Lists all server os-updates.

```
stackit beta server os-update list [flags]
```

### Examples

```
  List all os-updates for a server with ID "xxx"
  $ stackit beta server os-update list --server-id xxx

  List all os-updates for a server with ID "xxx" in JSON format
  $ stackit beta server os-update list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server os-update list"
      --limit int          Maximum number of entries to list
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

* [stackit beta server os-update](./stackit_beta_server_os-update.md)	 - Provides functionality for managed server updates

