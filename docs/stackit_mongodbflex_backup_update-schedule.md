## stackit mongodbflex backup update-schedule

Updates the backup schedule and retention policy for a MongoDB Flex instance

### Synopsis

Updates the backup schedule and retention policy for a MongoDB Flex instance.
The current backup schedule and retention policy can be seen in the output of the "stackit mongodbflex backup schedule" command.
The backup schedule is defined in the cron scheduling system format e.g. '0 0 * * *'.
See below for more detail on the retention policy options.

```
stackit mongodbflex backup update-schedule [flags]
```

### Examples

```
  Update the backup schedule of a MongoDB Flex instance with ID "xxx"
  $ stackit mongodbflex backup update-schedule --instance-id xxx --schedule '6 6 * * *'

  Update the retention days for backups of a MongoDB Flex instance with ID "xxx" to 5 days
  $ stackit mongodbflex backup update-schedule --instance-id xxx --store-for-days 5
```

### Options

```
  -h, --help                               Help for "stackit mongodbflex backup update-schedule"
      --instance-id string                 Instance ID
      --schedule string                    Backup schedule, in the cron scheduling system format e.g. '0 0 * * *'
      --store-daily-backup-days int        Number of days to retain daily backups. Should be less than or equal to the number of days of the selected weekly or monthly value.
      --store-for-days int                 Number of days to retain backups. Should be less than or equal to the value of the daily backup.
      --store-monthly-backups-months int   Number of months to retain monthly backups
      --store-weekly-backup-weeks int      Number of weeks to retain weekly backups. Should be less than or equal to the number of weeks of the selected monthly value.
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

* [stackit mongodbflex backup](./stackit_mongodbflex_backup.md)	 - Provides functionality for MongoDB Flex instance backups

