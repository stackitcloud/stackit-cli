## stackit beta logs instance list

Lists Logs instances

### Synopsis

Lists Logs instances within the project.

```
stackit beta logs instance list [flags]
```

### Examples

```
  List all Logs instances
  $ stackit beta logs instance list

  List the first 10 Logs instances
  $ stackit beta logs instance list --limit=10
```

### Options

```
  -h, --help        Help for "stackit beta logs instance list"
      --limit int   Limit the output to the first n elements
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

* [stackit beta logs instance](./stackit_beta_logs_instance.md)	 - Provides functionality for Logs instances

