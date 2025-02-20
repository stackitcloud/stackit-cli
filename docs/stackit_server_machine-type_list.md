## stackit server machine-type list

Get list of all machine types available in a project

### Synopsis

Get list of all machine types available in a project.

```
stackit server machine-type list [flags]
```

### Examples

```
  Get list of all machine types
  $ stackit server machine-type list

  Get list of all machine types in JSON format
  $ stackit server machine-type list --output-format json

  List the first 10 machine types
  $ stackit server machine-type list --limit=10
```

### Options

```
  -h, --help        Help for "stackit server machine-type list"
      --limit int   Limit the output to the first n elements
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

* [stackit server machine-type](./stackit_server_machine-type.md)	 - Provides functionality for server machine types available inside a project

