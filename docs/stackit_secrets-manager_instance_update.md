## stackit secrets-manager instance update

Updates a Secrets Manager instance

### Synopsis

Updates a Secrets Manager instance.

```
stackit secrets-manager instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the range of IPs allowed to access a Secrets Manager instance with ID "xxx"
  $ stackit secrets-manager instance update xxx --acl 1.2.3.0/24

  Update the KMS key settings of a Secrets Manager instance with ID "xxx"
  $ stackit secrets-manager instance update xxx --kms-key-id key-id --kms-keyring-id keyring-id --kms-key-version 1 --kms-service-account-email my-service-account-1234567@sa.stackit.cloud
```

### Options

```
      --acl strings                        List of IP networks in CIDR notation which are allowed to access this instance (default [])
  -h, --help                               Help for "stackit secrets-manager instance update"
      --kms-key-id string                  ID of the KMS key to use for encryption
      --kms-key-version int                Version of the KMS key
      --kms-keyring-id string              ID of the KMS key ring
      --kms-service-account-email string   Service account email for KMS access
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

