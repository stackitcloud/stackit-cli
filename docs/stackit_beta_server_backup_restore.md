## stackit beta server backup restore

Restores a Server Backup.

### Synopsis

Restores a Server Backup. Operation always is async.

```
stackit beta server backup restore BACKUP_ID [flags]
```

### Examples

```
  Restore a Server Backup with ID "xxx" for server "zzz"
  $ stackit beta server backup restore xxx --server-id=zzz

  Restore a Server Backup with ID "xxx" for server "zzz" and start the server afterwards
  $ stackit beta server backup restore xxx --server-id=zzz --start-server-after-restore
```

### Options

```
  -h, --help                         Help for "stackit beta server backup restore"
  -s, --server-id string             Server ID
  -u, --start-server-after-restore   Should the server start after the backup restoring.
  -i, --volume-ids strings           Backup volume IDs, as comma separated UUID values. (default [])
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

* [stackit beta server backup](./stackit_beta_server_backup.md)	 - Provides functionality for Server Backup

