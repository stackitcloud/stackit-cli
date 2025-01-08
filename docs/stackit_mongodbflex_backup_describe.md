## stackit mongodbflex backup describe

Shows details of a backup for a MongoDB Flex instance

### Synopsis

Shows details of a backup for a MongoDB Flex instance.

```
stackit mongodbflex backup describe BACKUP_ID [flags]
```

### Examples

```
  Get details of a backup with ID "xxx" for a MongoDB Flex instance with ID "yyy"
  $ stackit mongodbflex backup describe xxx --instance-id yyy

  Get details of a backup with ID "xxx" for a MongoDB Flex instance with ID "yyy" in JSON format
  $ stackit mongodbflex backup describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex backup describe"
      --instance-id string   Instance ID
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

* [stackit mongodbflex backup](./stackit_mongodbflex_backup.md)	 - Provides functionality for MongoDB Flex instance backups

