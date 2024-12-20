## stackit beta key-pair describe

Describes a key pair

### Synopsis

Describes a key pair.

```
stackit beta key-pair describe KEY_PAIR_NAME [flags]
```

### Examples

```
  Get details about a key pair with name "KEY_PAIR_NAME"
  $ stackit beta key-pair describe KEY_PAIR_NAME

  Get only the SSH public key of a key pair with name "KEY_PAIR_NAME"
  $ stackit beta key-pair describe KEY_PAIR_NAME --public-key
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

* [stackit beta key-pair](./stackit_beta_key-pair.md)	 - Provides functionality for SSH key pairs

