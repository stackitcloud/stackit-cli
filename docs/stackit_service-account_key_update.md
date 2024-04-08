## stackit service-account key update

Updates a service account key

### Synopsis

Updates a service account key.
You can temporarily activate or deactivate the key and/or update its date of expiration.

```
stackit service-account key update KEY_ID [flags]
```

### Examples

```
  Temporarily deactivate a key with ID "xxx" of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --deactivate

  Activate a key of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --activate

  Update the expiration date of a key of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key update xxx --email my-service-account-1234567@sa.stackit.cloud --expires-in-days 30
```

### Options

```
      --activate              If set, activates the service account key
      --deactivate            If set, temporarily deactivates the service account key
  -e, --email string          Service account email
      --expires-in-days int   Number of days until expiration
  -h, --help                  Help for "stackit service-account key update"
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

