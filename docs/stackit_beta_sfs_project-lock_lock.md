## stackit beta sfs project-lock lock

Enables lock for a project

### Synopsis

Enables lock for a project. Necessary for immutable snapshots and to prevent accidental deletion of resources.

```
stackit beta sfs project-lock lock [flags]
```

### Examples

```
  Enable lock for project
  $ stackit beta sfs project-lock lock
```

### Options

```
  -h, --help   Help for "stackit beta sfs project-lock lock"
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

* [stackit beta sfs project-lock](./stackit_beta_sfs_project-lock.md)	 - Provides functionality for SFS project locks

