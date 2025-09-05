## stackit beta kms key import

Import a KMS Key Version

### Synopsis

Import a new version to the given KMS key.

```
stackit beta kms key import [flags]
```

### Examples

```
  Import a new version for the given KMS Key "my-key"
  $ stakit beta kms key import --key-ring "my-keyring-id" --key "my-key-id" --wrapped-key "base64-encoded-wrapped-key-material" --wrapping-key-id "my-wrapping-key-id"
```

### Options

```
  -h, --help                     Help for "stackit beta kms key import"
      --key string               ID of the KMS Key
      --key-ring string          ID of the KMS Key Ring
      --wrapped-key string       The wrapped key material that has to be imported. Encoded in base64
      --wrapping-key-id string   he unique id of the wrapping key the key material has been wrapped with
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

- [stackit beta kms key](./stackit_beta_kms_key.md) - Manage KMS Keys
