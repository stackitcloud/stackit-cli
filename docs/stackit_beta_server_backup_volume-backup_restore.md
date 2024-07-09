## stackit beta server backup volume-backup restore

Restore a Server Volume Backup to a volume.

### Synopsis

Restore a Server Volume Backup to a volume. Operation always is async.

```
stackit beta server backup volume-backup restore VOLUME_BACKUP_ID [flags]
```

### Examples

```
  Restore a Server Volume Backup with ID "xxx" for server "zzz" and backup "bbb" to volume "rrr"
  $ stackit beta server backup volume-backup restore xxx --server-id=zzz --backup-id=bbb --restore-volume-id=rrr
```

### Options

```
  -b, --backup-id string           Backup ID
  -h, --help                       Help for "stackit beta server backup volume-backup restore"
  -r, --restore-volume-id string   Restore Volume ID
  -s, --server-id string           Server ID
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

* [stackit beta server backup volume-backup](./stackit_beta_server_backup_volume-backup.md)	 - Provides functionality for Server Backup Volume Backups

