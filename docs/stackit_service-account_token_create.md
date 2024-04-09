## stackit service-account token create

Creates an access token for a service account

### Synopsis

Creates an access token for a service account.
The access token can be then used for API calls (where enabled).
The token is only displayed upon creation, and it will not be recoverable later.

```
stackit service-account token create [flags]
```

### Examples

```
  Create an access token for the service account with email "my-service-account-1234567@sa.stackit.cloud" with a default time to live
  $ stackit service-account token create --email my-service-account-1234567@sa.stackit.cloud

  Create an access token for the service account with email "my-service-account-1234567@sa.stackit.cloud" with a time to live of 10 days
  $ stackit service-account token create --email my-service-account-1234567@sa.stackit.cloud --ttl-days 10
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account token create"
      --ttl-days int   How long (in days) the new access token is valid (default 90)
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

* [stackit service-account token](./stackit_service-account_token.md)	 - Provides functionality regarding service account tokens

