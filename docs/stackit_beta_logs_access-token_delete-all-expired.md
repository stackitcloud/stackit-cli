## stackit beta logs access-token delete-all-expired

Deletes all expired log access token

### Synopsis

Deletes all expired log access token.

```
stackit beta logs access-token delete-all-expired [flags]
```

### Examples

```
  Delete all expired access tokens in instance "xxx"
  $ stackit logs access-token delete --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit beta logs access-token delete-all-expired"
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

* [stackit beta logs access-token](./stackit_beta_logs_access-token.md)	 - Provides functionality for Logs access-tokens

