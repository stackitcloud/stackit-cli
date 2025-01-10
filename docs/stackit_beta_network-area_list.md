## stackit beta network-area list

Lists all STACKIT Network Areas (SNA) of an organization

### Synopsis

Lists all STACKIT Network Areas (SNA) of an organization.

```
stackit beta network-area list [flags]
```

### Examples

```
  Lists all network areas of organization "xxx"
  $ stackit beta network-area list --organization-id xxx

  Lists all network areas of organization "xxx" in JSON format
  $ stackit beta network-area list --organization-id xxx --output-format json

  Lists up to 10 network areas of organization "xxx"
  $ stackit beta network-area list --organization-id xxx --limit 10
```

### Options

```
  -h, --help                     Help for "stackit beta network-area list"
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

* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for STACKIT Network Area (SNA)

