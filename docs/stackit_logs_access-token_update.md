## stackit logs access-token update

Updates a access token

### Synopsis

Updates a access token.

```
stackit logs access-token update ACCESS_TOKEN_ID [flags]
```

### Examples

```
  Update access token with ID "xxx" with new name "access-token-1"
  $ stackit logs access-token update xxx --instance-id yyy --display-name access-token-1

  Update access token with ID "xxx" with new description "Access token for Service XY"
  $ stackit logs access-token update xxx --instance-id yyy --description "Access token for Service XY"
```

### Options

```
      --description string    Description of the access token
      --display-name string   Display name for the access token
  -h, --help                  Help for "stackit logs access-token update"
      --instance-id string    ID of the Logs instance
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

