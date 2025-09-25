## stackit beta kms version list

List all key versions

### Synopsis

List all versions of a given key.

```
stackit beta kms version list [flags]
```

### Examples

```
  List all key versions for the key "my-key-id" inside the key ring "my-key-ring-id"
  $ stackit beta kms version list --key "my-key-id" --key-ring "my-key-ring-id"

  List all key versions in JSON format
  $ stackit beta kms version list --key "my-key-id" --key-ring "my-key-ring-id" -o json
```

### Options

```
  -h, --help              Help for "stackit beta kms version list"
      --key string        ID of the key
      --key-ring string   ID of the KMS key ring
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

* [stackit beta kms version](./stackit_beta_kms_version.md)	 - Manage KMS key versions

