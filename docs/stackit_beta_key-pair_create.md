## stackit beta key-pair create

Creates a key pair

### Synopsis

Creates a key pair.

```
stackit beta key-pair create [flags]
```

### Examples

```
  Create a new key pair with public-key "ssh-rsa xxx"
  $ stackit beta key-pair create --public-key `ssh-rsa xxx`

  Create a new key pair with public-key from file "/Users/username/.ssh/id_rsa.pub"
  $ stackit beta key-pair create --public-key `@/Users/username/.ssh/id_rsa.pub`

  Create a new key pair with name "KEY_PAIR_NAME" and public-key "ssh-rsa yyy"
  $ stackit beta key-pair create --name KEY_PAIR_NAME --public-key `ssh-rsa yyy`

  Create a new key pair with public-key "ssh-rsa xxx" and labels "key=value,key1=value1"
  $ stackit beta key-pair create --public-key `ssh-rsa xxx` --labels key=value,key1=value1
```

### Options

```
  -h, --help                    Help for "stackit beta key-pair create"
      --labels stringToString   Labels are key-value string pairs which can be attached to a key pair. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --name string             key pair name
      --public-key string       Public key to be imported (format: ssh-rsa|ssh-ed25519)
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

