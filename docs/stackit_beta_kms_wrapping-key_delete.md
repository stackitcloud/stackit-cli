## stackit beta kms wrapping-key delete

Deletes a KMS wrapping key

### Synopsis

Deletes a KMS wrapping key inside a specific key ring.

```
stackit beta kms wrapping-key delete WRAPPING_KEY_ID [flags]
```

### Examples

```
  Delete a KMS wrapping key "MY_WRAPPING_KEY_ID" inside the key ring "MY_KEYRING_ID"
  $ stackit beta kms wrapping-key delete "MY_WRAPPING_KEY_ID" --keyring-id "MY_KEYRING_ID"
```

### Options

```
  -h, --help                Help for "stackit beta kms wrapping-key delete"
      --keyring-id string   ID of the KMS key ring where the wrapping key is stored
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

* [stackit beta kms wrapping-key](./stackit_beta_kms_wrapping-key.md)	 - Manage KMS wrapping keys

