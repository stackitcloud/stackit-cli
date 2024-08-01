## stackit beta network-area network-range delete

Deletes a network range in a STACKIT Network Area (SNA)

### Synopsis

Deletes a network range in a STACKIT Network Area (SNA).

```
stackit beta network-area network-range delete [flags]
```

### Examples

```
  Delete network range with id "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"
  $ stackit beta network-area network-range delete xxx --network-area-id yyy --organization-id zzz
```

### Options

```
  -h, --help                     Help for "stackit beta network-area network-range delete"
      --network-area-id string   STACKIT Network Area (SNA) ID
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

