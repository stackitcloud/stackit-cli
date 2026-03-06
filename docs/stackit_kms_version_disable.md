## stackit kms version disable

Disable a key version

### Synopsis

Disable the given key version.

```
stackit kms version disable VERSION_NUMBER [flags]
```

### Examples

```
  Disable key version "42" for the key "my-key-id" inside the key ring "my-keyring-id"
  $ stackit kms version disable 42 --key-id "my-key-id" --keyring-id "my-keyring-id"
```

### Options

```
  -h, --help                Help for "stackit kms version disable"
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

* [stackit kms version](./stackit_kms_version.md)	 - Manage KMS key versions

