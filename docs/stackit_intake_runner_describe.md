## stackit intake runner describe

Shows details of an Intake Runner

### Synopsis

Shows details of an Intake Runner.

```
stackit intake runner describe RUNNER_ID [flags]
```

### Examples

```
  Get details of an Intake Runner with ID "xxx"
  $ stackit intake runner describe xxx

  Get details of an Intake Runner with ID "xxx" in JSON format
  $ stackit intake runner describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit intake runner describe"
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

