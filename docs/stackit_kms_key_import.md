## stackit kms key import

Import a KMS key

### Synopsis

After encrypting the secret with the wrapping key’s public key and Base64-encoding it, import it as a new version of the specified KMS key.

```
stackit kms key import KEY_ID [flags]
```

### Examples

```
  Import a new version for the given KMS key "MY_KEY_ID" from literal value
  $ stackit kms key import "MY_KEY_ID" --keyring-id "my-keyring-id" --wrapped-key "BASE64_VALUE" --wrapping-key-id "MY_WRAPPING_KEY_ID"

  Import from a file
  $ stackit kms key import "MY_KEY_ID" --keyring-id "my-keyring-id" --wrapped-key "@path/to/wrapped.key.b64" --wrapping-key-id "MY_WRAPPING_KEY_ID"
```

### Options

```
  -h, --help                     Help for "stackit kms key import"
      --keyring-id string        ID of the KMS key ring
      --wrapped-key string       The wrapped key material to be imported. Base64-encoded. Pass the value directly or a file path (e.g. @path/to/wrapped.key.b64)
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

* [stackit kms key](./stackit_kms_key.md)	 - Manage KMS keys

