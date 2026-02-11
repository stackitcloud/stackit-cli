## stackit logs instance delete

Deletes the given Logs instance

### Synopsis

Deletes the given Logs instance.

```
stackit logs instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a Logs instance with ID "xxx"
  $ stackit logs instance delete "xxx"
```

### Options

```
  -h, --help   Help for "stackit logs instance delete"
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

* [stackit logs instance](./stackit_logs_instance.md)	 - Provides functionality for Logs instances

