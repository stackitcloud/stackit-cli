## stackit service-account key create

Creates a service account key

### Synopsis

Creates a service account key.
You can generate an RSA keypair and provide the public key.
If you do not provide a public key, the service will generate a new key-pair and the private key is included in the response. You won't be able to retrieve it later.

```
stackit service-account key create [flags]
```

### Examples

```
  Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud with no expiration date"
  $ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud

  Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud" expiring in 10 days
  $ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud --expires-in-days 10

  Create a key for the service account with email "my-service-account-1234567@sa.stackit.cloud" and provide the public key in a .pem file"
  $ stackit service-account key create --email my-service-account-1234567@sa.stackit.cloud --public-key @./public.pem
```

### Options

```
  -e, --email string          Service account email
      --expires-in-days int   Number of days until expiration. When omitted, the key is valid until deleted
  -h, --help                  Help for "stackit service-account key create"
      --public-key string     Public key of the user generated RSA 2048 key-pair. Must be in x509 format. Can be a string or path to the .pem file, if prefixed with "@"
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

* [stackit service-account key](./stackit_service-account_key.md)	 - Provides functionality regarding service account keys

