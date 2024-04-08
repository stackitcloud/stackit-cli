## stackit service-account get-jwks

Shows the JWKS for a service account

### Synopsis

Shows the JSON Web Key set (JWKS) for a service account. Only JSON output is supported.

```
stackit service-account get-jwks EMAIL [flags]
```

### Examples

```
  Get JWKS for the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account get-jwks my-service-account-1234567@sa.stackit.cloud
```

### Options

```
  -h, --help   Help for "stackit service-account get-jwks"
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

