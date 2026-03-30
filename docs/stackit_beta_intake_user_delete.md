## stackit beta intake user delete

Deletes an Intake User

### Synopsis

Deletes an Intake User.

```
stackit beta intake user delete USER_ID [flags]
```

### Examples

```
  Delete an Intake User with ID "xxx" for Intake "yyy"
  $ stackit beta intake user delete xxx --intake-id yyy
```

### Options

```
  -h, --help               Help for "stackit beta intake user delete"
      --intake-id string   Intake ID
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

