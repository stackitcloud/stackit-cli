## stackit beta sfs share create

Creates a share

### Synopsis

Creates a share.

```
stackit beta sfs share create [flags]
```

### Examples

```
  Create a share in a resource pool with ID "xxx", name "yyy" and no space hard limit
  $ stackit beta sfs share create --resource-pool-id xxx --name yyy --hard-limit 0

  Create a share in a resource pool with ID "xxx", name "yyy" and export policy with name "zzz"
  $ stackit beta sfs share create --resource-pool-id xxx --name yyy --export-policy-name zzz --hard-limit 0
```

### Options

```
      --export-policy-name string   The export policy the share is assigned to
      --hard-limit int              The space hard limit for the share
  -h, --help                        Help for "stackit beta sfs share create"
      --name string                 Share name
      --resource-pool-id string     The resource pool the share is assigned to
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

