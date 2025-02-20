## stackit server backup describe

Shows details of a Server Backup

### Synopsis

Shows details of a Server Backup.

```
stackit server backup describe BACKUP_ID [flags]
```

### Examples

```
  Get details of a Server Backup with id "my-backup-id"
  $ stackit server backup describe my-backup-id

  Get details of a Server Backup with id "my-backup-id" in JSON format
  $ stackit server backup describe my-backup-id --output-format json
```

### Options

```
  -h, --help               Help for "stackit server backup describe"
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

* [stackit server backup](./stackit_server_backup.md)	 - Provides functionality for server backups

