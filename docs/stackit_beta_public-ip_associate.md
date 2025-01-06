## stackit beta public-ip associate

Associates a Public IP with a network interface or a virtual IP

### Synopsis

Associates a Public IP with a network interface or a virtual IP.

```
stackit beta public-ip associate PUBLIC_IP_ID [flags]
```

### Examples

```
  Associate public IP with ID "xxx" to a resource (network interface or virtual IP) with ID "yyy"
  $ stackit beta public-ip associate xxx --associated-resource-id yyy
```

### Options

```
      --associated-resource-id string   Associates the public IP with a network interface or virtual IP (ID)
  -h, --help                            Help for "stackit beta public-ip associate"
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

* [stackit beta public-ip](./stackit_beta_public-ip.md)	 - Provides functionality for public IPs

