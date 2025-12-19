## stackit beta sfs export-policy delete

Deletes a export policy

### Synopsis

Deletes a export policy.

```
stackit beta sfs export-policy delete EXPORT_POLICY_ID [flags]
```

### Examples

```
  Delete a export policy with ID "xxx"
  $ stackit beta sfs export-policy delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta sfs export-policy delete"
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

* [stackit beta sfs export-policy](./stackit_beta_sfs_export-policy.md)	 - Provides functionality for SFS export policies

