## stackit observability instance list

Lists all Observability instances

### Synopsis

Lists all Observability instances.

```
stackit observability instance list [flags]
```

### Examples

```
  List all Observability instances
  $ stackit observability instance list

  List all Observability instances in JSON format
  $ stackit observability instance list --output-format json

  List up to 10 Observability instances
  $ stackit observability instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit observability instance list"
      --limit int   Maximum number of entries to list
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

* [stackit observability instance](./stackit_observability_instance.md)	 - Provides functionality for Observability instances

