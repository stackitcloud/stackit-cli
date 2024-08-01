## stackit beta network-area network-range create

Creates a network range in a STACKIT Network Area (SNA)

### Synopsis

Creates a network range in a STACKIT Network Area (SNA).

```
stackit beta network-area network-range create [flags]
```

### Examples

```
  Create a network range in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area network-range create --network-area-id xxx --organization-id yyy --network-range "1.1.1.0/24"
```

### Options

```
  -h, --help                     Help for "stackit beta network-area network-range create"
      --network-area-id string   STACKIT Network Area (SNA) ID
      --network-range string     Network range to create in CIDR notation
      --organization-id string   Organization ID
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

* [stackit beta network-area network-range](./stackit_beta_network-area_network-range.md)	 - Provides functionality for network ranges in STACKIT Network Areas

