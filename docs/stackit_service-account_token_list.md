## stackit service-account token list

List access tokens of a service account

### Synopsis

List access tokens of a service account.
Only the metadata about the access tokens is shown, and not the tokens themselves.
Access tokens (including revoked tokens) are returned until they are expired.

```
stackit service-account token list [flags]
```

### Examples

```
  List all access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud

  List all access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud" in JSON format
  $ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud --output-format json

  List up to 10 access tokens of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account token list --email my-service-account-1234567@sa.stackit.cloud --limit 10
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account token list"
      --limit int      Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit service-account token](./stackit_service-account_token.md)	 - Provides functionality regarding service account tokens

