## stackit secrets-manager instance create

Creates a Secrets Manager instance

### Synopsis

Creates a Secrets Manager instance.

```
stackit secrets-manager instance create [flags]
```

### Examples

```
  Create a Secrets Manager instance with name "my-instance"
  $ stackit secrets-manager instance create --name my-instance

  Create a Secrets Manager instance with name "my-instance" and specify IP range which is allowed to access it
  $ stackit secrets-manager instance create --name my-instance --acl 1.2.3.0/24

  Create a Secrets Manager instance with name "my-instance" and configure KMS key options
  $ stackit secrets-manager instance create --name my-instance --kms-key-id key-id --kms-keyring-id keyring-id --kms-key-version 1 --kms-service-account-email my-service-account-1234567@sa.stackit.cloud
```

### Options

```
      --acl strings                        List of IP networks in CIDR notation which are allowed to access this instance (default [])
  -h, --help                               Help for "stackit secrets-manager instance create"
      --kms-key-id string                  ID of the KMS key to use for encryption
      --kms-key-version int                Version of the KMS key
      --kms-keyring-id string              ID of the KMS key ring
      --kms-service-account-email string   Service account email for KMS access
  -n, --name string                        Instance name
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit secrets-manager instance](./stackit_secrets-manager_instance.md)	 - Provides functionality for Secrets Manager instances

