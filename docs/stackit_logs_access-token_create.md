## stackit logs access-token create

Creates a log access token

### Synopsis

Creates a log access token.

```
stackit logs access-token create [flags]
```

### Examples

```
  Create a access token with the display name "access-token-1" for the instance "xxx" with read and write permissions
  $ stackit logs access-token create --display-name access-token-1 --instance-id xxx --permissions read,write

  Create a write only access token with a description
  $ stackit logs access-token create --display-name access-token-2 --instance-id xxx --permissions write --description "Access token for service"

  Create a read only access token which expires in 30 days
  $ stackit logs access-token create --display-name access-token-3 --instance-id xxx --permissions read --lifetime 30
```

### Options

```
      --description string    Description of the access token
      --display-name string   Display name for the access token
  -h, --help                  Help for "stackit logs access-token create"
      --instance-id string    ID of the logs instance
      --lifetime int          Lifetime of the access token in days [1 - 180]
      --permissions strings   Permissions of the access token ["read" "write"]
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

* [stackit logs access-token](./stackit_logs_access-token.md)	 - Provides functionality for Logs access-tokens

