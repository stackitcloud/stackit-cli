## stackit beta server backup schedule list

Lists all server backup schedules

### Synopsis

Lists all server backup schedules.

```
stackit beta server backup schedule list [flags]
```

### Examples

```
  List all backup schedules for a server with ID "xxx"
  $ stackit beta server backup schedule list --server-id xxx

  List all backup schedules for a server with ID "xxx" in JSON format
  $ stackit beta server backup schedule list --server-id xxx --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server backup schedule list"
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

* [stackit beta server backup schedule](./stackit_beta_server_backup_schedule.md)	 - Provides functionality for Server Backup Schedule

