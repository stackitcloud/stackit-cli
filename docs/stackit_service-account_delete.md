## stackit service-account delete

Deletes a service account

### Synopsis

Deletes a service account.

```
stackit service-account delete EMAIL [flags]
```

### Examples

```
  Delete a service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account delete my-service-account-1234567@sa.stackit.cloud
```

### Options

```
  -h, --help   Help for "stackit service-account delete"
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

* [stackit service-account](./stackit_service-account.md)	 - Provides functionality for service accounts

