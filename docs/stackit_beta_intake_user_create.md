## stackit beta intake user create

Creates a new Intake User

### Synopsis

Creates a new Intake User for a specific Intake.

```
stackit beta intake user create [flags]
```

### Examples

```
  Create a new Intake User with required parameters
  $ stackit beta intake user create --display-name intake-user --intake-id xxx --password "SuperSafepass123\!"

  Create a new Intake User for the dead-letter queue with labels
  $ stackit beta intake user create --display-name dlq-user --intake-id xxx --password "SuperSafepass123\!" --type dead-letter --labels "env=prod"
```

### Options

```
      --description string      Description
      --display-name string     Display name
  -h, --help                    Help for "stackit beta intake user create"
      --intake-id string        The UUID of the Intake to associate the user with
      --labels stringToString   Labels in key=value format, separated by commas (default [])
      --password string         Password for the user. Must contain lower, upper, number, and special characters (min 12 chars)
      --type string             Type of user. One of 'intake' (default) or 'dead-letter' (default "intake")
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

