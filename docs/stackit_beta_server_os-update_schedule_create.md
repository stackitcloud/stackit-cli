## stackit beta server os-update schedule create

Creates a Server os-update Schedule

### Synopsis

Creates a Server os-update Schedule.

```
stackit beta server os-update schedule create [flags]
```

### Examples

```
  Create a Server os-update Schedule with name "myschedule"
  $ stackit beta server os-update schedule create --server-id xxx --name=myschedule

  Create a Server os-update Schedule with name "myschedule" and maintenance window for 14 o'clock
  $ stackit beta server os-update schedule create --server-id xxx --name=myschedule --maintenance-window=14

  Create a Server os-update Schedule with name "myschedule" and a rrule (os-update happening every week on Monday, at 22:30, for 30 weeks in a row, starting on 2025-01-01)
  $ stackit beta server os-update schedule create --server-id xxx --name=myschedule --rrule='DTSTART;TZID=Etc/UTC:20250101T223000 RRULE:FREQ=WEEKLY;INTERVAL=1;BYDAY=MO;COUNT=30'
```

### Options

```
  -e, --enabled                  Is the server os-update schedule enabled (default true)
  -h, --help                     Help for "stackit beta server os-update schedule create"
  -d, --maintenance-window int   os-update maintenance window (in hours, 1-24) (default 1)
  -n, --name string              os-update schedule name
  -r, --rrule string             os-update RRULE (recurrence rule) (default "DTSTART;TZID=Etc/UTC:20200803T023000 RRULE:FREQ=DAILY;INTERVAL=1")
  -s, --server-id string         Server ID
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

* [stackit beta server os-update schedule](./stackit_beta_server_os-update_schedule.md)	 - Provides functionality for Server os-update Schedule

