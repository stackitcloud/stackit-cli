## stackit beta kms keyring describe

Describe a KMS key ring

### Synopsis

Describe a KMS key ring

```
stackit beta kms keyring describe KEYRING_ID [flags]
```

### Examples

```
  Describe a KMS key ring with ID xxx
  $ stackit beta kms keyring describe xxx
```

### Options

```
  -h, --help   Help for "stackit beta kms keyring describe"
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

* [stackit beta kms keyring](./stackit_beta_kms_keyring.md)	 - Manage KMS key rings

