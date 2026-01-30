## stackit beta logs access-token list

Lists all access tokens of a project

### Synopsis

Lists all access tokens of a project.

```
stackit beta logs access-token list [flags]
```

### Examples

```
  Lists all access tokens of the instance "xxx"
  $ stackit logs access-token list --instance-id xxx

  Lists all access tokens in JSON format
  $ stackit logs access-token list --instance-id xxx --output-format json

  Lists up to 10 access-token
  $ stackit logs access-token list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit beta logs access-token list"
      --instance-id string   ID of the logs instance
      --limit int            Maximum number of entries to list
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

