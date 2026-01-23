## stackit beta edge-cloud instance list

Lists edge instances

### Synopsis

Lists STACKIT Edge Cloud (STEC) instances of a project.

```
stackit beta edge-cloud instance list [flags]
```

### Examples

```
  Lists all edge instances of a given project
  $ stackit beta edge-cloud instance list

  Lists all edge instances of a given project and limits the output to two instances
  $ stackit beta edge-cloud instance list --limit 2
```

### Options

```
  -h, --help        Help for "stackit beta edge-cloud instance list"
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

* [stackit beta edge-cloud instance](./stackit_beta_edge-cloud_instance.md)	 - Provides functionality for edge instances.

