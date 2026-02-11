## stackit logs access-token delete

Deletes a logs access token

### Synopsis

Deletes a logs access token.

```
stackit logs access-token delete ACCESS_TOKEN_ID [flags]
```

### Examples

```
  Delete access token with ID "xxx" in instance "yyy"
  $ stackit logs access-token delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit logs access-token delete"
      --instance-id string   ID of the logs instance
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

