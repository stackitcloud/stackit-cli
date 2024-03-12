## stackit postgresflex backups update-schedule

Updates backup schedule for a specific PostgreSQL Flex instance

### Synopsis

Updates backup schedule for a specific PostgreSQL Flex instance.

```
stackit postgresflex backups update-schedule INSTANCE_ID [flags]
```

### Examples

```
  Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex backups update-schedule xxx --backup-schedule '6 6 * * *'
```

### Options

```
      --backup-schedule string   Backup schedule
  -h, --help                     Help for "stackit postgresflex backups update-schedule"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex backups](./stackit_postgresflex_backups.md)	 - Provides functionality for PostgreSQL Flex instance backups

