## stackit beta network-interface update

Updates a network interface

### Synopsis

Updates a network interface.

```
stackit beta network-interface update NIC_ID [flags]
```

### Examples

```
  Updates a network interface with nic id "xxx" and network-id "yyy" to new allowed addresses "1.1.1.1,8.8.8.8,9.9.9.9" and new labels "key=value,key2=value2"
  $ stackit beta network-interface update xxx --network-id yyy --allowed-addresses "1.1.1.1,8.8.8.8,9.9.9.9" --labels key=value,key2=value2

  Updates a network interface with nic id "xxx" and network-id "yyy" with new name "nic-name-new"
  $ stackit beta network-interface update xxx --network-id yyy --name nic-name-new

  Updates a network interface with nic id "xxx" and network-id "yyy" to include the security group "zzz"
  $ stackit beta network-interface update xxx --network-id yyy --security-groups zzz
```

### Options

```
      --allowed-addresses strings   List of allowed IPs
  -h, --help                        Help for "stackit beta network-interface update"
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
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta network-interface](./stackit_beta_network-interface.md)	 - Provides functionality for network interfaces

