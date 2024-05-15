## stackit logme plans

Lists all LogMe service plans

### Synopsis

Lists all LogMe service plans.

```
stackit logme plans [flags]
```

### Examples

```
  List all LogMe service plans
  $ stackit logme plans

  List all LogMe service plans in JSON format
  $ stackit logme plans --output-format json

  List up to 10 LogMe service plans
  $ stackit logme plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit logme plans"
      --limit int   Maximum number of entries to list
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

* [stackit logme](./stackit_logme.md)	 - Provides functionality for LogMe

