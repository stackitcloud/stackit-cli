## stackit beta edge-cloud instance describe

Describes an edge instance

### Synopsis

Describes a STACKIT Edge Cloud (STEC) instance.

```
stackit beta edge-cloud instance describe [flags]
```

### Examples

```
  Describe an edge instance with id "xxx"
  $ stackit beta edge-cloud instance describe --id <ID>

  Describe an edge instance with name "xxx"
  $ stackit beta edge-cloud instance describe --name <NAME>
```

### Options

```
  -h, --help          Help for "stackit beta edge-cloud instance describe"
  -i, --id string     The project-unique identifier of this instance.
  -n, --name string   The displayed name to distinguish multiple instances.
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

