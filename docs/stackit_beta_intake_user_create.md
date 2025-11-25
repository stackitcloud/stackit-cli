## stackit beta intake user create

Creates a new Intake User

### Synopsis

Creates a new Intake User, providing secure access credentials for applications to connect to a data stream.

```
stackit beta intake user create [flags]
```

### Examples

```
  Create a new Intake User with a display name and password for a specific Intake
  $ stackit beta intake user create --intake-id xxx --display-name my-intake-user --password "my-secret-password"

  Create a new dead-letter queue user with a description and labels
  $ stackit beta intake user create --intake-id xxx --display-name my-dlq-reader --password "another-secret" --type "dead-letter" --description "User for reading undelivered messages" --labels "owner=team-alpha,scope=dlq"
```

### Options

```
      --description string      Description
      --display-name string     Display name
  -h, --help                    Help for "stackit beta intake user create"
      --intake-id string        ID of the Intake to which the user belongs
      --labels stringToString   Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2" (default [])
      --password string         User password
      --type string             Type of user, 'intake' for writing to the stream or 'dead-letter' for reading from the dead-letter queue
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

