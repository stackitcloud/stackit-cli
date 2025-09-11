## stackit beta kms key create

Creates a KMS key

### Synopsis

Creates a KMS key.

```
stackit beta kms key create [flags]
```

### Examples

```
  Create a Symmetric KMS key
  $ stackit beta kms key create --key-ring "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "my-key-name" --purpose "symmetric_encrypt_decrypt"

  Create a Message Authentication KMS key
  $ stackit beta kms key create --key-ring "my-keyring-id" --algorithm "hmac_sha512" --name "my-key-name" --purpose "message_authentication_code"
```

### Options

```
      --algorithm string     En-/Decryption / signing algorithm
      --description string   Optional description of the key
  -h, --help                 Help for "stackit beta kms key create"
      --import-only          States whether versions can be created or only imported
      --key-ring string      ID of the KMS key ring
      --name string          The display name to distinguish multiple keys
      --purpose string       Purpose of the key. Enum: 'symmetric_encrypt_decrypt', 'asymmetric_encrypt_decrypt', 'message_authentication_code', 'asymmetric_sign_verify' 
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

* [stackit beta kms key](./stackit_beta_kms_key.md)	 - Manage KMS Keys

