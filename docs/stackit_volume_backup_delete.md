## stackit volume backup delete

Deletes a backup

### Synopsis

Deletes a backup by its ID.

```
stackit volume backup delete BACKUP_ID [flags]
```

### Examples

```
  Delete a backup
  $ stackit volume backup delete xxx-xxx-xxx

  Delete a backup and wait for deletion to be completed
  $ stackit volume backup delete xxx-xxx-xxx --async=false
```

### Options

```
  -h, --help   Help for "stackit volume backup delete"
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

* [stackit volume backup](./stackit_volume_backup.md)	 - Provides functionality for volume backups

