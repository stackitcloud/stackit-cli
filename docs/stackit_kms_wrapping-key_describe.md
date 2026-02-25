## stackit kms wrapping-key describe

Describe a KMS wrapping key

### Synopsis

Describe a KMS wrapping key

```
stackit kms wrapping-key describe WRAPPING_KEY_ID [flags]
```

### Examples

```
  Describe a KMS wrapping key with ID xxx of keyring yyy
  $ stackit kms wrappingkey describe xxx --keyring-id yyy
```

### Options

```
  -h, --help                Help for "stackit kms wrapping-key describe"
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

* [stackit kms wrapping-key](./stackit_kms_wrapping-key.md)	 - Manage KMS wrapping keys

