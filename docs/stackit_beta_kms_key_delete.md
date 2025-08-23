## stackit beta kms key delete

Deletes a KMS Key

### Synopsis

Deletes a KMS Key inside a specific Key Ring.

```
stackit beta kms key delete [flags]
```

### Examples

```
  Delete a KMS Key "my-key-id" inside the Key Ring "my-key-ring-id"
  $ stackit beta kms keyring delete --key-ring "my-key-ring-id" --key "my-key-id"
```

### Options

```
  -h, --help              Help for "stackit beta kms key delete"
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

