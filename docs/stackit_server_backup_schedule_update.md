## stackit server backup schedule update

Updates a Server Backup Schedule

### Synopsis

Updates a Server Backup Schedule.

```
stackit server backup schedule update SCHEDULE_ID [flags]
```

### Examples

```
  Update the retention period of the backup schedule "zzz" of server "xxx"
  $ stackit server backup schedule update zzz --server-id=xxx --backup-retention-period=20

  Update the backup name of the backup schedule "zzz" of server "xxx"
  $ stackit server backup schedule update zzz --server-id=xxx --backup-name=newname
```

### Options

```
  -b, --backup-name string            Backup name
  -d, --backup-retention-period int   Backup retention period (in days) (default 14)
  -n, --backup-schedule-name string   Backup schedule name
  -i, --backup-volume-ids strings     Backup volume IDs, as comma separated UUID values. (default [])
  -e, --enabled                       Is the server backup schedule enabled (default true)
  -h, --help                          Help for "stackit server backup schedule update"
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

