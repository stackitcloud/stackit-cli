## stackit beta key-pair describe

Describe a keypair

### Synopsis

Describe a keypair.

```
stackit beta key-pair describe [flags]
```

### Examples

```
  Get details about a keypair named "KEYPAIR_NAME"
  $ stackit beta keypair describe KEYPAIR_NAME

  Get only the SSH public key of a keypair with the name "KEYPAIR_NAME"
  $ stackit beta keypair describe KEYPAIR_NAME --public-key
```

### Options

```
  -h, --help         Help for "stackit beta key-pair describe"
      --public-key   Show only the public key
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

