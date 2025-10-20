## stackit beta kms key import

Import a KMS key

### Synopsis

Import a new version to the given KMS key.

```
stackit beta kms key import KEY_ID [flags]
```

### Examples

```
  Import a new version for the given KMS key "my-key-id"
  $ stackit beta kms key import "my-key-id" --keyring-id "my-keyring-id" --wrapped-key "base64-encoded-wrapped-key-material" --wrapping-key-id "my-wrapping-key-id"
```

### Options

```
  -h, --help                     Help for "stackit beta kms key import"
      --keyring-id string        ID of the KMS key ring
      --wrapped-key string       The wrapped key material that has to be imported. Encoded in base64
      --wrapping-key-id string   The unique id of the wrapping key the key material has been wrapped with
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

