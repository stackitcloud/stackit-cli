## stackit beta network-interface create

Creates a network interface

### Synopsis

Creates a network interface.

```
stackit beta network-interface create [flags]
```

### Examples

```
  Create a network interface for network with ID "xxx"
  $ stackit beta network-interface create --network-id xxx

  Create a network interface with allowed addresses, labels, a name, security groups and nic security enabled for network with ID "xxx"
  $ stackit beta network-interface create --network-id xxx --allowed-addresses "1.1.1.1,8.8.8.8,9.9.9.9" --labels key=value,key2=value2 --name NAME --security-groups "UUID1,UUID2" --nic-security
```

### Options

```
      --allowed-addresses strings   List of allowed IPs
  -h, --help                        Help for "stackit beta network-interface create"
  -i, --ipv4 string                 IPv4 address
  -s, --ipv6 string                 IPv6 address
      --labels stringToString       Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string                 Network interface name
      --network-id string           Network ID
  -b, --nic-security                If this is set to false, then no security groups will apply to this network interface. (default true)
      --security-groups strings     List of security groups
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

* [stackit beta network-interface](./stackit_beta_network-interface.md)	 - Provides functionality for network interfaces

