## stackit logme instance delete

Deletes a LogMe instance

### Synopsis

Deletes a LogMe instance.

```
stackit logme instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a LogMe instance with ID "xxx"
  $ stackit logme instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit logme instance delete"
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

* [stackit logme instance](./stackit_logme_instance.md)	 - Provides functionality for LogMe instances

