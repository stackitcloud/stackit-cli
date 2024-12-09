## stackit beta key-pair create

Create a keypair

### Synopsis

Create a keypair.

```
stackit beta key-pair create [flags]
```

### Examples

```
  Create a new key-pair with public-key "ssh-rsa xxx"
  $ stackit beta key-pair create --public-key ssh-rsa xxx

  Create a new key-pair with public-key from file "/Users/username/.ssh/id_rsa.pub"
  $ stackit beta key-pair create --public-key @/Users/username/.ssh/id_rsa.pub

  Create a new key-pair with name "KEYPAIR_NAME" and public-key "ssh-rsa yyy"
  $ stackit beta key-pair create --name KEYPAIR_NAME --public-key ssh-rsa yyy

  Create a new key-pair with public-key "ssh-rsa xxx" and labels "key=value,key1=value1"
  $ stackit beta key-pair create --public-key ssh-rsa xxx --labels key=value,key1=value1
```

### Options

```
  -h, --help                    Help for "stackit beta key-pair create"
      --labels stringToString   Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --name string             Name of the key which will be created
      --public-key string       Public key which should be add (format: ssh-rsa|sha-ed25519)
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

