## stackit beta sfs share update

Updates a share

### Synopsis

Updates a share.

```
stackit beta sfs share update SHARE_ID [flags]
```

### Examples

```
  Update share with ID "xxx" with new export-policy-name "yyy" in resource-pool "zzz"
  $ stackit beta sfs share update xxx --export-policy-name yyy --resource-pool-id zzz

  Update share with ID "xxx" with new space hard limit "50" in resource-pool "yyy"
  $ stackit beta sfs share update xxx --hard-limit 50 --resource-pool-id yyy
```

### Options

```
      --export-policy-name string   The export policy the share is assigned to
      --hard-limit int              The space hard limit for the share
  -h, --help                        Help for "stackit beta sfs share update"
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

