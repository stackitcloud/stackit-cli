## stackit beta kms wrappingkey create

Creates a KMS Wrapping Key

### Synopsis

Creates a KMS Wrapping Key.

```
stackit beta kms wrappingkey create [flags]
```

### Examples

```
  Create a Symmetric KMS Wrapping Key
  $ stakit beta kms wrappingkey create --key-ring "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "my-wrapping-key-name" --purpose "wrap_symmetric_key"

  Create an Asymmetric KMS Wrapping Key with a description
  $ stakit beta kms wrappingkey create --key-ring "my-keyring-id" --algorithm "hmac_sha256" --name "my-wrapping-key-name" --description "my-description" --purpose "wrap_asymmetric_key"
```

### Options

```
      --algorithm string     En-/Decryption algorithm
      --backend string       The backend that is responsible for maintaining this wrapping key (default "software")
      --description string   Optinal description of the Wrapping Key
  -h, --help                 Help for "stackit beta kms wrappingkey create"
      --key-ring string      ID of the KMS Key Ring
      --name string          The display name to distinguish multiple wrapping keys
      --purpose string       Purpose of the Wrapping Key. Enum: 'wrap_symmetric_key', 'wrap_asymmetric_key' 
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

