## stackit observability plans

Lists all Observability service plans

### Synopsis

Lists all Observability service plans.

```
stackit observability plans [flags]
```

### Examples

```
  List all Observability service plans
  $ stackit observability plans

  List all Observability service plans in JSON format
  $ stackit observability plans --output-format json

  List up to 10 Observability service plans
  $ stackit observability plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit observability plans"
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

* [stackit observability](./stackit_observability.md)	 - Provides functionality for Observability

