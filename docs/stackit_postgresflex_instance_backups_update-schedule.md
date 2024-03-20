## stackit postgresflex instance backups update-schedule

Updates backup schedule for a specific PostgreSQL Flex instance

### Synopsis

Updates backup schedule for a specific PostgreSQL Flex instance.

```
stackit postgresflex instance backups update-schedule [flags]
```

### Examples

```
  Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex instance backups update-schedule --instance-id xxx --backup-schedule '6 6 * * *'
```

### Options

```
      --backup-schedule string   Backup schedule, in the cron scheduling system format e.g. '0 0 * * *'
  -h, --help                     Help for "stackit postgresflex instance backups update-schedule"
      --instance-id string       Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex instance backups](./stackit_postgresflex_instance_backups.md)	 - Provides functionality for PostgreSQL Flex instance backups

