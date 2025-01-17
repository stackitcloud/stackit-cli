## stackit beta server os-update schedule list

Lists all server os-update schedules

### Synopsis

Lists all server os-update schedules.

```
stackit beta server os-update schedule list [flags]
```

### Examples

```
  List all os-update schedules for a server with ID "xxx"
  $ stackit beta server os-update schedule list --server-id xxx

  List all os-update schedules for a server with ID "xxx" in JSON format
  $ stackit beta server os-update schedule list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server os-update schedule list"
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

* [stackit beta server os-update schedule](./stackit_beta_server_os-update_schedule.md)	 - Provides functionality for Server os-update Schedule

