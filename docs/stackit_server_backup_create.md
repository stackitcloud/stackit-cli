## stackit server backup create

Creates a Server Backup.

### Synopsis

Creates a Server Backup. Operation always is async.

```
stackit server backup create [flags]
```

### Examples

```
  Create a Server Backup with name "mybackup"
  $ stackit server backup create --server-id xxx --name=mybackup

  Create a Server Backup with name "mybackup" and retention period of 5 days
  $ stackit server backup create --server-id xxx --name=mybackup --retention-period=5
```

### Options

```
  -h, --help                   Help for "stackit server backup create"
  -b, --name string            Backup name
  -d, --retention-period int   Backup retention period (in days) (default 14)
  -s, --server-id string       Server ID
  -i, --volume-ids strings     Backup volume IDs, as comma separated UUID values. (default [])
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

