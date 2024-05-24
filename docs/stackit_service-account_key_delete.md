## stackit service-account key delete

Deletes a service account key

### Synopsis

Deletes a service account key.

```
stackit service-account key delete KEY_ID [flags]
```

### Examples

```
  Delete a key with ID "xxx" belonging to the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account key delete  xxx --email my-service-account-1234567@sa.stackit.cloud
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account key delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit service-account key](./stackit_service-account_key.md)	 - Provides functionality for service account keys

