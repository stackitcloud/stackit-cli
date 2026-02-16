## stackit beta intake user list

Lists all Intake Users

### Synopsis

Lists all Intake Users for a specific Intake.

```
stackit beta intake user list [flags]
```

### Examples

```
  List all users for an Intake
  $ stackit beta intake user list --intake-id xxx

  List all users for an Intake in JSON format
  $ stackit beta intake user list --intake-id xxx --output-format json

  List up to 5 users for an Intake
  $ stackit beta intake user list --intake-id xxx --limit 5
```

### Options

```
  -h, --help               Help for "stackit beta intake user list"
      --intake-id string   Intake ID
      --limit int          Maximum number of entries to list
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

* [stackit beta intake user](./stackit_beta_intake_user.md)	 - Provides functionality for Intake Users

