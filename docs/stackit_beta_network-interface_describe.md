## stackit beta network-interface describe

Describes a network interface

### Synopsis

Describes a network interface.

```
stackit beta network-interface describe [flags]
```

### Examples

```
  Describes network interface with nic id "xxx" and network ID "yyy"
  $ stackit beta network-interface describe xxx --network-id yyy

  Describes network interface with nic id "xxx" and network ID "yyy" in JSON format
  $ stackit beta network-interface describe xxx --network-id yyy --output-format json

  Describes network interface with nic id "xxx" and network ID "yyy" in yaml format
  $ stackit beta network-interface describe xxx --network-id yyy --output-format yaml
```

### Options

```
  -h, --help                Help for "stackit beta network-interface describe"
      --network-id string   Network ID
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

* [stackit beta network-interface](./stackit_beta_network-interface.md)	 - Provides functionality for Network Interface

