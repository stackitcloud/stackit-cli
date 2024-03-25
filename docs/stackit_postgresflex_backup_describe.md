## stackit postgresflex backup describe

Shows details of a backup for a specific PostgreSQL Flex instance

### Synopsis

Shows details of a backup for a specific PostgreSQL Flex instance.

```
stackit postgresflex backup describe BACKUP_ID [flags]
```

### Examples

```
  Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy"
  $ stackit postgresflex backup describe xxx --instance-id yyy

  Get details of a backup with ID "xxx" for a PostgreSQL Flex instance with ID "yyy" in a table format
  $ stackit postgresflex backup describe xxx --instance-id yyy --output-format pretty
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
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex backup](./stackit_postgresflex_backup.md)	 - Provides functionality for PostgreSQL Flex instance backups

