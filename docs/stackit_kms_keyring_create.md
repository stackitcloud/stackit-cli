## stackit kms keyring create

Creates a KMS key ring

### Synopsis

Creates a KMS key ring.

```
stackit kms keyring create [flags]
```

### Examples

```
  Create a KMS key ring with name "my-keyring"
  $ stackit kms keyring create --name my-keyring

  Create a KMS key ring with a description
  $ stackit kms keyring create --name my-keyring --description my-description

  Create a KMS key ring and print the result as YAML
  $ stackit kms keyring create --name my-keyring -o yaml
```

### Options

```
      --description string   Optional description of the key ring
  -h, --help                 Help for "stackit kms keyring create"
      --name string          Name of the KMS key ring
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

