## stackit server backup schedule create

Creates a Server Backup Schedule

### Synopsis

Creates a Server Backup Schedule.

```
stackit server backup schedule create [flags]
```

### Examples

```
  Create a Server Backup Schedule with name "myschedule" and backup name "mybackup"
  $ stackit server backup schedule create --server-id xxx --backup-name=mybackup --backup-schedule-name=myschedule

  Create a Server Backup Schedule with name "myschedule", backup name "mybackup" and retention period of 5 days
  $ stackit server backup schedule create --server-id xxx --backup-name=mybackup --backup-schedule-name=myschedule --backup-retention-period=5
```

### Options

```
  -b, --backup-name string            Backup name
  -d, --backup-retention-period int   Backup retention period (in days) (default 14)
  -n, --backup-schedule-name string   Backup schedule name
  -i, --backup-volume-ids strings     Backup volume IDs, as comma separated UUID values. (default [])
  -e, --enabled                       Is the server backup schedule enabled (default true)
  -h, --help                          Help for "stackit server backup schedule create"
  -r, --rrule string                  Backup RRULE (recurrence rule) (default "DTSTART;TZID=Europe/Sofia:20200803T023000 RRULE:FREQ=DAILY;INTERVAL=1")
  -s, --server-id string              Server ID
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

* [stackit server backup schedule](./stackit_server_backup_schedule.md)	 - Provides functionality for Server Backup Schedule

