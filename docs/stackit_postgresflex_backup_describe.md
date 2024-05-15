## stackit postgresflex backup describe

Shows details of a backup for a PostgreSQL Flex instance

### Synopsis

Shows details of a backup for a PostgreSQL Flex instance.

```
stackit postgresflex backup describe BACKUP_ID [flags]
```

### Examples

```
  Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy"
  $ stackit postgresflex backup describe xxx --instance-id yyy

  Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy" in JSON format
  $ stackit postgresflex backup describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit postgresflex backup describe"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit postgresflex backup](./stackit_postgresflex_backup.md)	 - Provides functionality for PostgreSQL Flex instance backups

