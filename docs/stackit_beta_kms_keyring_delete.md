## stackit beta kms keyring delete

Deletes a KMS Keyring

### Synopsis

Deletes a KMS Keyring.

```
stackit beta kms keyring delete KEYRING_ID [flags]
```

### Examples

```
  Delete a KMS Keyring with ID "xxx"
  $ stackit beta kms keyring delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta kms keyring delete"
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

* [stackit beta kms keyring](./stackit_beta_kms_keyring.md)	 - Manage KMS Keyrings

