## stackit beta sfs resource-pool list

Lists all SFS resource pools

### Synopsis

Lists all SFS resource pools.

```
stackit beta sfs resource-pool list [flags]
```

### Examples

```
  List all SFS resource pools
  $ stackit beta sfs resource-pool list

  List all SFS resource pools for another region than the default one
  $ stackit beta sfs resource-pool list --region eu01

  List up to 10 SFS resource pools
  $ stackit beta sfs resource-pool list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta sfs resource-pool list"
      --limit int   Maximum number of entries to list
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

* [stackit beta sfs resource-pool](./stackit_beta_sfs_resource-pool.md)	 - Provides functionality for SFS resource pools

