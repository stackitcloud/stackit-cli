## stackit kms key create

Creates a KMS key

### Synopsis

Creates a KMS key.

```
stackit kms key create [flags]
```

### Examples

```
  Create a symmetric AES key (AES-256) with the name "symm-aes-gcm" under the key ring "my-keyring-id"
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "aes_256_gcm" --name "symm-aes-gcm" --purpose "symmetric_encrypt_decrypt" --protection "software"

  Create an asymmetric RSA encryption key (RSA-2048)
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "prod-orders-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software"

  Create a message authentication key (HMAC-SHA512)
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "hmac_sha512" --name "api-mac-key" --purpose "message_authentication_code" --protection "software"

  Create an ECDSA P-256 key for signing & verification
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "ecdsa_p256_sha256" --name "signing-ecdsa-p256" --purpose "asymmetric_sign_verify" --protection "software"

  Create an import-only key (versions must be imported)
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "ext-managed-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software" --import-only

  Create a key and print the result as YAML
  $ stackit kms key create --keyring-id "my-keyring-id" --algorithm "rsa_2048_oaep_sha256" --name "yaml-output-rsa" --purpose "asymmetric_encrypt_decrypt" --protection "software" --output yaml
```

### Options

```
      --algorithm string     En-/Decryption / signing algorithm. Possible values: ["aes_256_gcm" "rsa_2048_oaep_sha256" "rsa_3072_oaep_sha256" "rsa_4096_oaep_sha256" "rsa_4096_oaep_sha512" "hmac_sha256" "hmac_sha384" "hmac_sha512" "ecdsa_p256_sha256" "ecdsa_p384_sha384" "ecdsa_p521_sha512"]
      --description string   Optional description of the key
  -h, --help                 Help for "stackit kms key create"
      --import-only          States whether versions can be created or only imported
      --keyring-id string    ID of the KMS key ring
      --name string          The display name to distinguish multiple keys
      --protection string    The underlying system that is responsible for protecting the key material. Possible values: ["symmetric_encrypt_decrypt" "asymmetric_encrypt_decrypt" "message_authentication_code" "asymmetric_sign_verify"]
      --purpose string       Purpose of the key. Possible values: ["symmetric_encrypt_decrypt" "asymmetric_encrypt_decrypt" "message_authentication_code" "asymmetric_sign_verify"]
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

