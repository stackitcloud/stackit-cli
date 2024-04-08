## stackit service-account key list

Lists all service account keys

### Synopsis

Lists all service account keys.

```
stackit service-account key list [flags]
```

### Examples

```
  List all keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud

  List all keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud" in JSON format
  $ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud --output-format json

  List up to 10 keys belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key list --email my-service-account-1234567@sa.stackit.cloud --limit 10
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account key list"
      --limit int      Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit service-account key](./stackit_service-account_key.md)	 - Provides functionality regarding service account keys

