## stackit server machine-type describe

Shows details of a server machine type

### Synopsis

Shows details of a server machine type.

```
stackit server machine-type describe MACHINE_TYPE [flags]
```

### Examples

```
  Show details of a server machine type with name "xxx"
  $ stackit server machine-type describe xxx

  Show details of a server machine type with name "xxx" in JSON format
  $ stackit server machine-type describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit server machine-type describe"
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

