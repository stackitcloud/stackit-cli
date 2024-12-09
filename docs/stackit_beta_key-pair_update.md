## stackit beta key-pair update

Update key pair

### Synopsis

Update key pair.

```
stackit beta key-pair update [flags]
```

### Examples

```
  Update the labels of a keypair KEYPAIR_NAME with "key=value,key1=value1"
  $ stackit beta key-pair update --labels key=value,key1=value1
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
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta key-pair](./stackit_beta_key-pair.md)	 - Provides functionality for Keypairs

