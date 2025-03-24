## stackit server os-update create

Creates a Server os-update.

### Synopsis

Creates a Server os-update. Operation always is async.

```
stackit server os-update create [flags]
```

### Examples

```
  Create a Server os-update with name "myupdate"
  $ stackit server os-update create --server-id xxx

  Create a Server os-update with name "myupdate" and maintenance window for 13 o'clock.
  $ stackit server os-update create --server-id xxx --maintenance-window=13
```

### Options

```
  -h, --help                     Help for "stackit server os-update create"
  -m, --maintenance-window int   Maintenance window (in hours, 1-24) (default 23)
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

* [stackit server os-update](./stackit_server_os-update.md)	 - Provides functionality for managed server updates

