## stackit volume backup list

Lists all backups

### Synopsis

Lists all backups in a project.

```
stackit volume backup list [flags]
```

### Examples

```
  List all backups
  $ stackit volume backup list

  List all backups in JSON format
  $ stackit volume backup list --output-format json

  List up to 10 backups
  $ stackit volume backup list --limit 10

  List backups with specific labels
  $ stackit volume backup list --label-selector key1=value1,key2=value2
```

### Options

```
  -h, --help                    Help for "stackit volume backup list"
      --label-selector string   Filter backups by labels
      --limit int               Maximum number of entries to list
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

