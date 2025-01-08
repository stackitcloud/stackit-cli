## stackit observability credentials list

Lists the usernames of all credentials for an Observability instance

### Synopsis

Lists the usernames of all credentials for an Observability instance.

```
stackit observability credentials list [flags]
```

### Examples

```
  List the usernames of all credentials for an Observability instance with ID "xxx"
  $ stackit observability credentials list --instance-id xxx

  List the usernames of all credentials for an Observability instance in JSON format
  $ stackit observability credentials list --instance-id xxx --output-format json

  List the usernames of up to 10 credentials for an Observability instance
  $ stackit observability credentials list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit observability credentials list"
      --instance-id string   Instance ID
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

* [stackit observability credentials](./stackit_observability_credentials.md)	 - Provides functionality for Observability credentials

