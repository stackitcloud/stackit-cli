## stackit logme instance list

Lists all LogMe instances

### Synopsis

Lists all LogMe instances.

```
stackit logme instance list [flags]
```

### Examples

```
  List all LogMe instances
  $ stackit logme instance list

  List all LogMe instances in JSON format
  $ stackit logme instance list --output-format json

  List up to 10 LogMe instances
  $ stackit logme instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit logme instance list"
      --limit int   Maximum number of entries to list
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

