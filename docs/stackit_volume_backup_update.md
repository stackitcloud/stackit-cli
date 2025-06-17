## stackit volume backup update

Updates a backup

### Synopsis

Updates a backup by its ID.

```
stackit volume backup update BACKUP_ID [flags]
```

### Examples

```
  Update the name of a backup with ID "xxx"
  $ stackit volume backup update xxx --name new-name

  Update the labels of a backup with ID "xxx"
  $ stackit volume backup update xxx --labels key1=value1,key2=value2
```

### Options

```
  -h, --help                    Help for "stackit volume backup update"
      --labels stringToString   Key-value string pairs as labels (default [])
      --name string             Name of the backup
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

