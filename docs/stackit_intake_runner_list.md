## stackit intake runner list

Lists all Intake Runners

### Synopsis

Lists all Intake Runners for the current project.

```
stackit intake runner list [flags]
```

### Examples

```
  List all Intake Runners
  $ stackit intake runner list

  List all Intake Runners in JSON format
  $ stackit intake runner list --output-format json

  List up to 5 Intake Runners
  $ stackit intake runner list --limit 5
```

### Options

```
  -h, --help        Help for "stackit intake runner list"
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

* [stackit intake runner](./stackit_intake_runner.md)	 - Provides functionality for Intake Runners

