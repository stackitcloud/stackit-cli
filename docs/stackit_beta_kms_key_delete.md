## stackit beta kms key delete

Deletes a KMS key

### Synopsis

Deletes a KMS key inside a specific key ring.

```
stackit beta kms key delete KEY_ID [flags]
```

### Examples

```
  Delete a KMS key "my-key-id" inside the key ring "my-key-ring-id"
  $ stackit beta kms key delete "my-key-id" --key-ring "my-key-ring-id"
```

### Options

```
  -h, --help                 Help for "stackit beta kms key delete"
      --key-ring-id string   ID of the KMS Key Ring where the Key is stored
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

* [stackit beta kms key](./stackit_beta_kms_key.md)	 - Manage KMS Keys

