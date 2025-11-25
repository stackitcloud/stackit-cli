## stackit beta intake user describe

Shows details of an Intake User

### Synopsis

Shows details of an Intake User.

```
stackit beta intake user describe USER_ID [flags]
```

### Examples

```
  Get details of an Intake User with ID "xxx" from an Intake with ID "yyy"
  $ stackit beta intake user describe xxx --intake-id yyy

  Get details of an Intake User in JSON format
  $ stackit beta intake user describe xxx --intake-id yyy --output-format json
```

### Options

```
  -h, --help               Help for "stackit beta intake user describe"
      --intake-id string   ID of the Intake to which the user belongs
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

