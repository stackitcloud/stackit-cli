## stackit beta edge-cloud plans list

Lists available edge service plans

### Synopsis

Lists available STACKIT Edge Cloud (STEC) service plans of a project

```
stackit beta edge-cloud plans list [flags]
```

### Examples

```
  Lists all edge plans for a given project
  $ stackit beta edge-cloud plan list

  Lists all edge plans for a given project and limits the output to two plans
  $ stackit beta edge-cloud plan list --limit 2
```

### Options

```
  -h, --help        Help for "stackit beta edge-cloud plans list"
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

* [stackit beta edge-cloud plans](./stackit_beta_edge-cloud_plans.md)	 - Provides functionality for edge service plans.

