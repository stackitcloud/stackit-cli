## stackit beta sfs share delete

Deletes a share

### Synopsis

Deletes a share.

```
stackit beta sfs share delete SHARE_ID [flags]
```

### Examples

```
  Delete a share with ID "xxx" from a resource pool with ID "yyy"
  $ stackit beta sfs share delete xxx --resource-pool-id yyy
```

### Options

```
  -h, --help                      Help for "stackit beta sfs share delete"
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

