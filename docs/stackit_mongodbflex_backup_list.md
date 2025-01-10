## stackit mongodbflex backup list

Lists all backups which are available for a MongoDB Flex instance

### Synopsis

Lists all backups which are available for a MongoDB Flex instance.

```
stackit mongodbflex backup list [flags]
```

### Examples

```
  List all backups of instance with ID "xxx"
  $ stackit mongodbflex backup list --instance-id xxx

  List all backups of instance with ID "xxx" in JSON format
  $ stackit mongodbflex backup list --instance-id xxx --output-format json

  List up to 10 backups of instance with ID "xxx"
  $ stackit mongodbflex backup list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex backup list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
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

