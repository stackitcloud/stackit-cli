## stackit observability instance delete

Deletes an Observability instance

### Synopsis

Deletes an Observability instance.

```
stackit observability instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete an Observability instance with ID "xxx"
  $ stackit Observability instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit observability instance delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit observability instance](./stackit_observability_instance.md)	 - Provides functionality for Observability instances

