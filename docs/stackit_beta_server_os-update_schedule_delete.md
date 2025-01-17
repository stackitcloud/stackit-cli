## stackit beta server os-update schedule delete

Deletes a Server os-update Schedule

### Synopsis

Deletes a Server os-update Schedule.

```
stackit beta server os-update schedule delete SCHEDULE_ID [flags]
```

### Examples

```
  Delete a Server os-update Schedule with ID "xxx" for server "zzz"
  $ stackit beta server os-update schedule delete xxx --server-id=zzz
```

### Options

```
  -h, --help               Help for "stackit beta server os-update schedule delete"
  -s, --server-id string   Server ID
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

