## stackit beta intake describe

Shows details of an Intake

### Synopsis

Shows details of an Intake.

```
stackit beta intake describe INTAKE_ID [flags]
```

### Examples

```
  Get details of an Intake with ID "xxx"
  $ stackit beta intake describe xxx

  Get details of an Intake with ID "xxx" in JSON format
  $ stackit beta intake describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta intake describe"
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

