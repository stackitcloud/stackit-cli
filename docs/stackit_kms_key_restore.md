## stackit kms key restore

Restore a key

### Synopsis

Restores the given key from deletion.

```
stackit kms key restore KEY_ID [flags]
```

### Examples

```
  Restore a KMS key "MY_KEY_ID" inside the key ring "my-keyring-id" that was scheduled for deletion.
  $ stackit kms key restore "MY_KEY_ID" --keyring-id "my-keyring-id"
```

### Options

```
  -h, --help                Help for "stackit kms key restore"
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

