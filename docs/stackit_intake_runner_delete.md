## stackit intake runner delete

Deletes an Intake Runner

### Synopsis

Deletes an Intake Runner.

```
stackit intake runner delete RUNNER_ID [flags]
```

### Examples

```
  Delete an Intake Runner with ID "xxx"
  $ stackit intake runner delete xxx
```

### Options

```
  -h, --help   Help for "stackit intake runner delete"
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

