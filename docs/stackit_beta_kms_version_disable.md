## stackit beta kms version disable

Disable a key version

### Synopsis

Disable the given key version.

```
stackit beta kms version disable [flags]
```

### Examples

```
  Disable key version "0" for the key "my-key-id" inside the key ring "my-key-ring-id"
  $ stackit beta kms version disable --key "my-key-id" --key-ring "my-key-ring-id" --version 0
```

### Options

```
  -h, --help              Help for "stackit beta kms version disable"
      --key string        ID of the rey
      --key-ring string   ID of the KMS key ring
      --version int       Version number of the key
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

