## stackit network-area list

Lists all STACKIT Network Areas (SNA) of an organization

### Synopsis

Lists all STACKIT Network Areas (SNA) of an organization.

```
stackit network-area list [flags]
```

### Examples

```
  Lists all network areas of organization "xxx"
  $ stackit network-area list --organization-id xxx

  Lists all network areas of organization "xxx" in JSON format
  $ stackit network-area list --organization-id xxx --output-format json

  Lists up to 10 network areas of organization "xxx"
  $ stackit network-area list --organization-id xxx --limit 10

  Lists all network areas of organization "xxx" which contains the label yyy
  $ stackit network-area list --organization-id xxx --label-selector yyy
```

### Options

```
  -h, --help                     Help for "stackit network-area list"
      --label-selector string    Filter by label
      --limit int                Maximum number of entries to list
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

