## stackit service-account key describe

Shows details of a service account key

### Synopsis

Shows details of a service account key. Only JSON output is supported.

```
stackit service-account key describe KEY_ID [flags]
```

### Examples

```
  Get details of a service account key with ID "xxx" belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key describe xxx --email my-service-account-1234567@sa.stackit.cloud
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account key describe"
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

* [stackit service-account key](./stackit_service-account_key.md)	 - Provides functionality for service account keys

