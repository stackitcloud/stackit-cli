## stackit beta public-ip disassociate

Disassociates a Public IP from a network interface or a virtual IP

### Synopsis

Disassociates a Public IP from a network interface or a virtual IP.

```
stackit beta public-ip disassociate [flags]
```

### Examples

```
  Disassociate public IP with ID "xxx" from a resource (network interface or virtual IP)
  $ stackit beta public-ip disassociate xxx
```

### Options

```
  -h, --help   Help for "stackit beta public-ip disassociate"
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

* [stackit beta public-ip](./stackit_beta_public-ip.md)	 - Provides functionality for Public IP

