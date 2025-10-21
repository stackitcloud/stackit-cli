## stackit beta kms key rotate

Rotate a key

### Synopsis

Rotates the given key.

```
stackit beta kms key rotate KEY_ID [flags]
```

### Examples

```
  Rotate a KMS key "MY_KEY_ID" and increase its version inside the key ring "MY_KEYRING_ID".
  $ stackit beta kms key rotate "MY_KEY_ID" --keyring-id "MY_KEYRING_ID"
```

### Options

```
  -h, --help                Help for "stackit beta kms key rotate"
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

* [stackit beta kms key](./stackit_beta_kms_key.md)	 - Manage KMS keys

