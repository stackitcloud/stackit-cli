## stackit beta sfs share list

Lists all shares of a resource pool

### Synopsis

Lists all shares of a resource pool.

```
stackit beta sfs share list [flags]
```

### Examples

```
  List all shares from resource pool with ID "xxx"
  $ stackit beta sfs export-policy list --resource-pool-id xxx

  List up to 10 shares from resource pool with ID "xxx"
  $ stackit beta sfs export-policy list --resource-pool-id xxx --limit 10
```

### Options

```
  -h, --help                      Help for "stackit beta sfs share list"
      --limit int                 Maximum number of entries to list
      --resource-pool-id string   The resource pool the share is assigned to
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

* [stackit beta sfs share](./stackit_beta_sfs_share.md)	 - Provides functionality for SFS shares

