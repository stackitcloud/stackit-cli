## stackit beta key-pair update

Updates a key pair

### Synopsis

Updates a key pair.

```
stackit beta key-pair update KEY_PAIR_NAME [flags]
```

### Examples

```
  Update the labels of a key pair with name "KEY_PAIR_NAME" with "key=value,key1=value1"
  $ stackit beta key-pair update KEY_PAIR_NAME --labels key=value,key1=value1
```

### Options

```
  -h, --help                    Help for "stackit beta key-pair update"
      --labels stringToString   Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...' (default [])
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

* [stackit beta key-pair](./stackit_beta_key-pair.md)	 - Provides functionality for SSH key pairs

