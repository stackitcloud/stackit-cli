## stackit logme credentials describe

Shows details of credentials of a LogMe instance

### Synopsis

Shows details of credentials of a LogMe instance. The password will be shown in plain text in the output.

```
stackit logme credentials describe CREDENTIALS_ID [flags]
```

### Examples

```
  Get details of credentials with ID "xxx" from instance with ID "yyy"
  $ stackit logme credentials describe xxx --instance-id yyy

  Get details of credentials with ID "xxx" from instance with ID "yyy" in JSON format
  $ stackit logme credentials describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit logme credentials describe"
      --instance-id string   Instance ID
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

* [stackit logme credentials](./stackit_logme_credentials.md)	 - Provides functionality for LogMe credentials

