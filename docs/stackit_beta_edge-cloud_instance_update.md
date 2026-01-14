## stackit beta edge-cloud instance update

Updates an edge instance

### Synopsis

Updates a STACKIT Edge Cloud (STEC) instance.

```
stackit beta edge-cloud instance update [flags]
```

### Examples

```
  Updates the description of an edge instance with id "xxx"
  $ stackit beta edge-cloud instance update --id "xxx" --description "yyy"

  Updates the plan of an edge instance with name "xxx"
  $ stackit beta edge-cloud instance update --name "xxx" --plan-id "yyy"

  Updates the description and plan of an edge instance with id "xxx"
  $ stackit beta edge-cloud instance update --id "xxx" --description "yyy" --plan-id "zzz"
```

### Options

```
  -d, --description string   A user chosen description to distinguish multiple instances.
  -h, --help                 Help for "stackit beta edge-cloud instance update"
  -i, --id string            The project-unique identifier of this instance.
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

