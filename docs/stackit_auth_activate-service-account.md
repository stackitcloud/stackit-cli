## stackit auth activate-service-account

Authenticates using a service account

### Synopsis

Authenticates to the CLI using service account credentials.
Subsequent commands will be authenticated using the service account credentials provided.
For more details on how to configure your service account, check our Authentication guide at https://github.com/stackitcloud/stackit-cli/blob/main/AUTHENTICATION.md.

```
stackit auth activate-service-account [flags]
```

### Examples

```
  Activate service account authentication in the STACKIT CLI using a service account key which includes the private key
  $ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json

  Activate service account authentication in the STACKIT CLI using the service account key and explicitly providing the private key in a PEM encoded file, which will take precedence over the one in the service account key
  $ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json --private-key-path path/to/private.key

  Activate service account authentication in the STACKIT CLI using the service account token
  $ stackit auth activate-service-account --service-account-token my-service-account-token
```

### Options

```
  -h, --help                              Help for "stackit auth activate-service-account"
      --jwks-custom-endpoint string       Custom endpoint for the jwks API, which is used to get the json web key sets (jwks) to validate tokens when the service-account authentication is activated
      --private-key-path string           RSA private key path. It takes precedence over the private key included in the service account key, if present
      --service-account-key-path string   Service account key path
      --service-account-token string      Service account long-lived access token
      --token-custom-endpoint string      Custom endpoint for the token API, which is used to request access tokens when the service-account authentication is activated
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

* [stackit auth](./stackit_auth.md)	 - Authenticates the STACKIT CLI

