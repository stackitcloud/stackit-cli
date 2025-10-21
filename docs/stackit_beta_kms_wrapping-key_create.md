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
  $ stackit beta kms wrapping-key create --keyring-id "MY_KEYRING_ID" --algorithm "rsa_2048_oaep_sha256" --name "my-wrapping-key-name" --purpose "wrap_asymmetric_key" --protection "software"

  Create an Asymmetric KMS wrapping key with a description
  $ stackit beta kms wrapping-key create --keyring-id "MY_KEYRING_ID" --algorithm "hmac_sha256" --name "my-wrapping-key-name" --description "my-description" --purpose "wrap_asymmetric_key" --protection "software"
```

### Options

```
      --algorithm string     En-/Decryption / signing algorithm. Possible values: ["rsa_2048_oaep_sha256" "rsa_3072_oaep_sha256" "rsa_4096_oaep_sha256" "rsa_4096_oaep_sha512" "rsa_2048_oaep_sha256_aes_256_key_wrap" "rsa_3072_oaep_sha256_aes_256_key_wrap" "rsa_4096_oaep_sha256_aes_256_key_wrap" "rsa_4096_oaep_sha512_aes_256_key_wrap"]
      --description string   Optional description of the wrapping key
  -h, --help                 Help for "stackit beta kms wrapping-key create"
      --keyring-id string    ID of the KMS key ring
      --name string          The display name to distinguish multiple wrapping keys
      --protection string    The underlying system that is responsible for protecting the wrapping key material. Possible values: ["wrap_symmetric_key" "wrap_asymmetric_key"]
      --purpose string       Purpose of the wrapping key. Possible values: ["wrap_symmetric_key" "wrap_asymmetric_key"]
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

