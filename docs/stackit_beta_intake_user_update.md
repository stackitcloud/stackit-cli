## stackit beta intake user update

Updates an Intake User

### Synopsis

Updates an Intake User. Only the specified fields are updated.

```
stackit beta intake user update USER_ID [flags]
```

### Examples

```
  Update the display name of an Intake User with ID "xxx"
  $ stackit beta intake user update xxx --intake-id yyy --display-name "new-user-name"

  Update the password of an Intake User
  $ stackit beta intake user update xxx --intake-id yyy --password "new-secret"
```

### Options

```
      --description string      Description
      --display-name string     Display name
  -h, --help                    Help for "stackit beta intake user update"
      --intake-id string        ID of the Intake
      --labels stringToString   Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2". (default [])
      --password string         User password
      --type string             Type of user, 'intake' for writing or 'dead-letter' for reading from the dead-letter queue
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

