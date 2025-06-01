## stackit volume backup create

Creates a backup from a specific source

### Synopsis

Creates a backup from a specific source (volume or snapshot).

```
stackit volume backup create [flags]
```

### Examples

```
  Create a backup from a volume
  $ stackit volume backup create --source-id xxx --source-type volume --project-id xxx

  Create a backup from a snapshot with a name
  $ stackit volume backup create --source-id xxx --source-type snapshot --name my-backup --project-id xxx

  Create a backup with labels
  $ stackit volume backup create --source-id xxx --source-type volume --labels key1=value1,key2=value2 --project-id xxx
```

### Options

```
  -h, --help                    Help for "stackit volume backup create"
      --labels stringToString   Key-value string pairs as labels (default [])
      --name string             Name of the backup
      --source-id string        ID of the source from which a backup should be created
      --source-type string      Source type of the backup (volume or snapshot)
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

