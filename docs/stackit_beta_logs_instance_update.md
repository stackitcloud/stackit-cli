## stackit beta logs instance update

Updates a Logs instance

### Synopsis

Updates a Logs instance.

```
stackit beta logs instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the display name of the Logs instance with ID "xxx"
  $ stackit beta logs instance update xxx --display-name new-name

  Update the retention time of the Logs instance with ID "xxx"
  $ stackit beta logs instance update xxx --retention-days 40

  Update the ACL of the Logs instance with ID "xxx"
  $ stackit beta logs instance update xxx --acl 1.2.3.0/24
```

### Options

```
      --acl strings           Access control list
      --description string    Description
      --display-name string   Display name
  -h, --help                  Help for "stackit beta logs instance update"
      --retention-days int    The days for how long the logs should be stored before being cleaned up
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

* [stackit beta logs instance](./stackit_beta_logs_instance.md)	 - Provides functionality for Logs instances

