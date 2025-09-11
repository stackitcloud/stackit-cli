## stackit beta kms keyring list

Lists all KMS key rings

### Synopsis

Lists all KMS key rings.

```
stackit beta kms keyring list [flags]
```

### Examples

```
  List all KMS key rings
  $ stackit beta kms keyring list

  List all KMS key rings in JSON format
  $ stackit beta kms keyring list --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta kms keyring list"
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

