## stackit beta server backup volume-backup delete

Deletes a Server Volume Backup.

### Synopsis

Deletes a Server Volume Backup. Operation always is async.

```
stackit beta server backup volume-backup delete VOLUME_BACKUP_ID [flags]
```

### Examples

```
  Delete a Server Volume Backup with ID "xxx" for server "zzz" and backup "bbb"
  $ stackit beta server backup volume-backup delete xxx --server-id=zzz --backup-id=bbb
```

### Options

```
  -b, --backup-id string   Backup ID
  -h, --help               Help for "stackit beta server backup volume-backup delete"
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

* [stackit beta server backup volume-backup](./stackit_beta_server_backup_volume-backup.md)	 - Provides functionality for Server Backup Volume Backups

