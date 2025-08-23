## stackit beta kms wrappingkey delete

Deletes a KMS Wrapping Key

### Synopsis

Deletes a KMS Wrapping Key inside a specific Key Ring.

```
stackit beta kms wrappingkey delete [flags]
```

### Examples

```
  Delete a KMS Wrapping Key "my-wrapping-key-id" inside the Key Ring "my-key-ring-id"
  $ stackit beta kms keyring delete --key-ring "my-key-ring-id" --wrapping-key "my-wrapping-key-id"
```

### Options

```
  -h, --help                  Help for "stackit beta kms wrappingkey delete"
      --key-ring string       ID of the KMS Key Ring where the Wrapping Key is stored
      --wrapping-key string   ID of the actual Wrapping Key
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

* [stackit beta kms wrappingkey](./stackit_beta_kms_wrappingkey.md)	 - Manage KMS Wrapping Keys

