## stackit beta server backup schedule delete

Deletes a Server Backup Schedule

### Synopsis

Deletes a Server Backup Schedule.

```
stackit beta server backup schedule delete SCHEDULE_ID [flags]
```

### Examples

```
  Delete a Server Backup Schedule with ID "xxx" for server "zzz"
  $ stackit beta server backup schedule delete xxx --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit beta server backup schedule delete"
  -s, --server-id string   Server ID
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

* [stackit beta server backup schedule](./stackit_beta_server_backup_schedule.md)	 - Provides functionality for Server Backup Schedule

