## stackit beta kms version restore

Restore a key version

### Synopsis

Restores the specified version of a key.

```
stackit beta kms version restore VERSION_NUMBER [flags]
```

### Examples

```
  Restore key version "42" for the key "MY_KEY_ID" inside the key ring "MY_KEYRING_ID"
  $ stackit beta kms version restore 42 --key-id "MY_KEY_ID" --keyring-id "MY_KEYRING_ID"
```

### Options

```
  -h, --help                Help for "stackit beta kms version restore"
      --key-id string       ID of the key
      --keyring-id string   ID of the KMS key ring
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

