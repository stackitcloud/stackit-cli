## stackit logme credentials list

Lists all credentials' IDs for a LogMe instance

### Synopsis

Lists all credentials' IDs for a LogMe instance.

```
stackit logme credentials list [flags]
```

### Examples

```
  List all credentials' IDs for a LogMe instance
  $ stackit logme credentials list --instance-id xxx

  List all credentials' IDs for a LogMe instance in JSON format
  $ stackit logme credentials list --instance-id xxx --output-format json

  List up to 10 credentials' IDs for a LogMe instance
  $ stackit logme credentials list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit logme credentials list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit logme credentials](./stackit_logme_credentials.md)	 - Provides functionality for LogMe credentials

