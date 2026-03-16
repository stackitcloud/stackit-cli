## stackit kms key list

List all KMS keys

### Synopsis

List all KMS keys inside a key ring.

```
stackit kms key list [flags]
```

### Examples

```
  List all KMS keys for the key ring "my-keyring-id"
  $ stackit kms key list --keyring-id "my-keyring-id"

  List all KMS keys in JSON format
  $ stackit kms key list --keyring-id "my-keyring-id" --output-format json
```

### Options

```
  -h, --help                Help for "stackit kms key list"
      --keyring-id string   ID of the KMS key ring where the key is stored
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

* [stackit kms key](./stackit_kms_key.md)	 - Manage KMS keys

