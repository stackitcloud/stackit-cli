## stackit service-account token revoke

Revokes an access token of a service account

### Synopsis

Revokes an access token of a service account.
The access token is instantly revoked, any following calls with the token will be unauthorized.
The token metadata is still stored until the expiration time.

```
stackit service-account token revoke TOKEN_ID [flags]
```

### Examples

```
  Revoke an access token with ID "xxx" of the service account with email "my-service-account-1234567@sa.stackit.cloud"
  $ stackit service-account token revoke xxx --email my-service-account-1234567@sa.stackit.cloud
```

### Options

```
  -e, --email string   Service account email
  -h, --help           Help for "stackit service-account token revoke"
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

