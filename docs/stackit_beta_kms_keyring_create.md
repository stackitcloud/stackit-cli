## stackit beta kms keyring create

Creates a KMS Key Ring

### Synopsis

Creates a KMS Key Ring.

```
stackit beta kms keyring create [flags]
```

### Examples

```
  Create a KMS key ring
  $ stakit beta kms keyring create --name my-keyring

  Create a KMS Key ring with a description
  $ stakit beta kms keyring create --name my-keyring --description my-description
```

### Options

```
      --description string   Optinal description of the Key Ring
  -h, --help                 Help for "stackit beta kms keyring create"
      --name string          Name of the KMS Key Ring
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

