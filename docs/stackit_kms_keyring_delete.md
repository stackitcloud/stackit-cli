## stackit kms keyring delete

Deletes a KMS key ring

### Synopsis

Deletes a KMS key ring.

```
stackit kms keyring delete KEYRING-ID [flags]
```

### Examples

```
  Delete a KMS key ring with ID "MY_KEYRING_ID"
  $ stackit kms keyring delete "MY_KEYRING_ID"
```

### Options

```
  -h, --help   Help for "stackit kms keyring delete"
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

* [stackit kms keyring](./stackit_kms_keyring.md)	 - Manage KMS key rings

