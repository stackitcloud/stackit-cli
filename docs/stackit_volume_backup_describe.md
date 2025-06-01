## stackit volume backup describe

Describes a backup

### Synopsis

Describes a backup by its ID.

```
stackit volume backup describe BACKUP_ID [flags]
```

### Examples

```
  Get details of a backup
  $ stackit volume backup describe xxx-xxx-xxx

  Get details of a backup in JSON format
  $ stackit volume backup describe xxx-xxx-xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit volume backup describe"
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

