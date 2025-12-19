## stackit beta sfs export-policy list

Lists all export policies of a project

### Synopsis

Lists all export policies of a project.

```
stackit beta sfs export-policy list [flags]
```

### Examples

```
  List all export policies
  $ stackit beta sfs export-policy list

  List up to 10 export policies
  $ stackit beta sfs export-policy list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta sfs export-policy list"
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

* [stackit beta sfs export-policy](./stackit_beta_sfs_export-policy.md)	 - Provides functionality for SFS export policies

