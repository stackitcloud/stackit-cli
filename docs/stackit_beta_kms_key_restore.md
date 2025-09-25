## stackit beta kms key restore

Resotre a key

### Synopsis

Restores the given key from being deleted.

```
stackit beta kms key restore [flags]
```

### Examples

```
  Restore a KMS key "my-key-id" inside the key ring "my-key-ring-id" that was scheduled for deletion.
  $ stackit beta kms keyring restore --key-ring "my-key-ring-id" --key "my-key-id"
```

### Options

```
  -h, --help              Help for "stackit beta kms key restore"
      --key string        ID of the actual Key
      --key-ring string   ID of the KMS Key Ring where the Key is stored
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

