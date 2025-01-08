## stackit beta server backup delete

Deletes a Server Backup.

### Synopsis

Deletes a Server Backup. Operation always is async.

```
stackit beta server backup delete BACKUP_ID [flags]
```

### Examples

```
  Delete a Server Backup with ID "xxx" for server "zzz"
  $ stackit beta server backup delete xxx --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit beta server backup delete"
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

* [stackit beta server backup](./stackit_beta_server_backup.md)	 - Provides functionality for server backups

