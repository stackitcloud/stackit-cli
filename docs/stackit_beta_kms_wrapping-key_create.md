## stackit beta kms wrapping-key create

Creates a KMS wrapping key

### Synopsis

Creates a KMS wrapping key.

```
stackit beta kms wrapping-key create [flags]
```

### Examples

```
  Create a Symmetric KMS wrapping key
  $ stackit beta kms wrappingkey create --key-ring "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "my-wrapping-key-name" --purpose "wrap_symmetric_key" --protection "software"

  Create an Asymmetric KMS wrapping key with a description
  $ stackit beta kms wrappingkey create --key-ring "my-keyring-id" --algorithm "hmac_sha256" --name "my-wrapping-key-name" --description "my-description" --purpose "wrap_asymmetric_key" --protection "software"
```

### Options

```
      --algorithm string     En-/Decryption algorithm
      --description string   Optional description of the wrapping key
  -h, --help                 Help for "stackit beta kms wrapping-key create"
      --key-ring-id string   ID of the KMS key ring
      --name string          The display name to distinguish multiple wrapping keys
      --protection string    Protection of the wrapping key. Value: 'software' 
      --purpose string       Purpose of the wrapping key. Enum: 'wrap_symmetric_key', 'wrap_asymmetric_key' 
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

