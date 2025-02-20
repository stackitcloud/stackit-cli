## stackit network-area delete

Deletes a STACKIT Network Area (SNA)

### Synopsis

Deletes a STACKIT Network Area (SNA) in an organization.
If the SNA is attached to any projects, the deletion will fail


```
stackit network-area delete AREA_ID [flags]
```

### Examples

```
  Delete network area with ID "xxx" in organization with ID "yyy"
  $ stackit network-area delete xxx --organization-id yyy
```

### Options

```
  -h, --help                     Help for "stackit network-area delete"
      --organization-id string   Organization ID
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

* [stackit network-area](./stackit_network-area.md)	 - Provides functionality for STACKIT Network Area (SNA)

