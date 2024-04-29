## stackit postgresflex backup update-schedule

Updates backup schedule for a PostgreSQL Flex instance

### Synopsis

Updates backup schedule for a PostgreSQL Flex instance. The current backup schedule can be seen in the output of the "stackit postgresflex instance describe" command.

```
stackit postgresflex backup update-schedule [flags]
```

### Examples

```
  Update the backup schedule of a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex backup update-schedule --instance-id xxx --schedule '6 6 * * *'
```

### Options

```
  -h, --help                 Help for "stackit postgresflex backup update-schedule"
      --instance-id string   Instance ID
      --schedule string      Backup schedule, in the cron scheduling system format e.g. '0 0 * * *'
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

* [stackit postgresflex backup](./stackit_postgresflex_backup.md)	 - Provides functionality for PostgreSQL Flex instance backups

