## stackit beta server backup schedule describe

Shows details of a Server Backup Schedule

### Synopsis

Shows details of a Server Backup Schedule.

```
stackit beta server backup schedule describe BACKUP_SCHEDULE_ID [flags]
```

### Examples

```
  Get details of a Server Backup Schedule with id "my-schedule-id"
  $ stackit beta server backup schedule describe my-schedule-id

  Get details of a Server Backup Schedule with id "my-schedule-id" in JSON format
  $ stackit beta server backup schedule describe my-schedule-id --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta server backup schedule describe"
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

