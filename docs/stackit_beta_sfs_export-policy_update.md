## stackit beta sfs export-policy update

Updates a export policy

### Synopsis

Updates a export policy.

```
stackit beta sfs export-policy update EXPORT_POLICY_ID [flags]
```

### Examples

```
  Update a export policy with ID "xxx" and with rules from file "./rules.json"
  $ stackit beta sfs export-policy update xxx --rules @./rules.json

  Update a export policy with ID "xxx" and remove the rules
  $ stackit beta sfs export-policy update XXX --remove-rules
```

### Options

```
  -h, --help           Help for "stackit beta sfs export-policy update"
      --remove-rules   Remove the export policy rules
      --rules string   Rules of the export policy
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

