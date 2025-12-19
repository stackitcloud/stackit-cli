## stackit beta sfs export-policy describe

Shows details of a export policy

### Synopsis

Shows details of a export policy.

```
stackit beta sfs export-policy describe EXPORT_POLICY_ID [flags]
```

### Examples

```
  Describe a export policy with ID "xxx"
  $ stackit beta sfs export-policy describe xxx
```

### Options

```
  -h, --help   Help for "stackit beta sfs export-policy describe"
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

