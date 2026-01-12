## stackit beta edge-cloud instance create

Creates an edge instance

### Synopsis

Creates a STACKIT Edge Cloud (STEC) instance. The instance will take a moment to become fully functional.

```
stackit beta edge-cloud instance create [flags]
```

### Examples

```
  Creates an edge instance with the name "xxx" and plan-id "yyy"
  $ stackit beta edge-cloud instance create --name "xxx" --plan-id "yyy"
```

### Options

```
  -d, --description string   A user chosen description to distinguish multiple instances.
  -h, --help                 Help for "stackit beta edge-cloud instance create"
  -n, --name string          The displayed name to distinguish multiple instances.
      --plan-id string       Service Plan configures the size of the Instance.
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

* [stackit beta edge-cloud instance](./stackit_beta_edge-cloud_instance.md)	 - Provides functionality for edge instances.

