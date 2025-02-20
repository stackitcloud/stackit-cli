## stackit network-interface describe

Describes a network interface

### Synopsis

Describes a network interface.

```
stackit network-interface describe NIC_ID [flags]
```

### Examples

```
  Describes network interface with nic id "xxx" and network ID "yyy"
  $ stackit network-interface describe xxx --network-id yyy

  Describes network interface with nic id "xxx" and network ID "yyy" in JSON format
  $ stackit network-interface describe xxx --network-id yyy --output-format json

  Describes network interface with nic id "xxx" and network ID "yyy" in yaml format
  $ stackit network-interface describe xxx --network-id yyy --output-format yaml
```

### Options

```
  -h, --help                Help for "stackit network-interface describe"
      --network-id string   Network ID
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

* [stackit network-interface](./stackit_network-interface.md)	 - Provides functionality for network interfaces

