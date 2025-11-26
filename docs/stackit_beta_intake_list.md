## stackit beta intake list

Lists all Intakes

### Synopsis

Lists all Intakes for the current project.

```
stackit beta intake list [flags]
```

### Examples

```
  List all Intakes
  $ stackit beta intake list

  List all Intakes in JSON format
  $ stackit beta intake list --output-format json

  List up to 5 Intakes
  $ stackit beta intake list --limit 5
```

### Options

```
  -h, --help        Help for "stackit beta intake list"
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

* [stackit beta intake](./stackit_beta_intake.md)	 - Provides functionality for intake

