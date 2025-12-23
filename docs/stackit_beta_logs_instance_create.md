## stackit beta logs instance create

Creates a Logs instance

### Synopsis

Creates a Logs instance.

```
stackit beta logs instance create [flags]
```

### Examples

```
  Create a Logs instance with name "my-instance" and retention time 10 days
  $ stackit beta logs instance create --display-name "my-instance" --retention-days 10

  Create a Logs instance with name "my-instance", retention time 10 days, and a description
  $ stackit beta logs instance create --display-name "my-instance" --retention-days 10 --description "Description of the instance"

  Create a Logs instance with name "my-instance", retention time 10 days, and restrict access to a specific range of IP addresses.
  $ stackit beta logs instance create --display-name "my-instance" --retention-days 10 --acl 1.2.3.0/24
```

### Options

```
      --acl strings           Access control list
      --description string    Description
      --display-name string   Display name
  -h, --help                  Help for "stackit beta logs instance create"
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

