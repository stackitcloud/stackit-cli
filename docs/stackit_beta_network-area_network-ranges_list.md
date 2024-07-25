## stackit beta network-area network-ranges list

Lists all network ranges in a STACKIT Network Area (SNA)

### Synopsis

Lists all network ranges in a STACKIT Network Area (SNA).

```
stackit beta network-area network-ranges list [flags]
```

### Examples

```
  Lists all network ranges in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area network-ranges list --network-area-id xxx --organization-id yyy

  Lists all network ranges in a STACKIT Network Area with ID "xxx" in organization with ID "yyy" in JSON format
  $ stackit beta network-area network-ranges list --network-area-id xxx --organization-id yyy --output-format json

  Lists up to 10 network ranges in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area network-ranges list --network-area-id xxx --organization-id yyy --limit 10
```

### Options

```
  -h, --help                     Help for "stackit beta network-area network-ranges list"
      --limit int                Maximum number of entries to list
      --network-area-id string   STACKIT Network Area ID
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

* [stackit beta network-area network-ranges](./stackit_beta_network-area_network-ranges.md)	 - Provides functionality for Network Area

