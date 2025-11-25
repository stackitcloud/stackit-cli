## stackit beta kms key describe

Describe a KMS key

### Synopsis

Describe a KMS key

```
stackit beta kms key describe KEY_ID [flags]
```

### Examples

```
  Describe a KMS key with ID xxx of keyring yyy
  $ stackit beta kms key describe xxx --keyring-id yyy
```

### Options

```
  -h, --help                Help for "stackit beta kms key describe"
      --keyring-id string   Key Ring ID
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

