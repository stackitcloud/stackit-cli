## stackit beta sfs export-policy create

Creates a export policy

### Synopsis

Creates a export policy.

```
stackit beta sfs export-policy create [flags]
```

### Examples

```
  Create a new export policy with name "EXPORT_POLICY_NAME"
  $ stackit beta sfs export-policy create --name EXPORT_POLICY_NAME

  Create a new export policy with name "EXPORT_POLICY_NAME" and rules from file "./rules.json"
  $ stackit beta sfs export-policy create --name EXPORT_POLICY_NAME --rules @./rules.json
```

### Options

```
  -h, --help           Help for "stackit beta sfs export-policy create"
      --name string    Export policy name
      --rules string   Rules of the export policy (format: json)
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

