## stackit mongodbflex backup restore-jobs

Lists all restore jobs which have been run for a MongoDB Flex instance

### Synopsis

Lists all restore jobs which have been run for a MongoDB Flex instance.

```
stackit mongodbflex backup restore-jobs [flags]
```

### Examples

```
  List all restore jobs of instance with ID "xxx"
  $ stackit mongodbflex backup restore-jobs --instance-id xxx

  List all restore jobs of instance with ID "xxx" in JSON format
  $ stackit mongodbflex backup restore-jobs --instance-id xxx --output-format json

  List up to 10 restore jobs of instance with ID "xxx"
  $ stackit mongodbflex backup restore-jobs --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex backup restore-jobs"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit mongodbflex backup](./stackit_mongodbflex_backup.md)	 - Provides functionality for MongoDB Flex instance backups

