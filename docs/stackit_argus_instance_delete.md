## stackit argus instance delete

Deletes an Argus instance

### Synopsis

Deletes an Argus instance.

```
stackit argus instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete an Argus instance with ID "xxx"
  $ stackit argus instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit argus instance delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus instance](./stackit_argus_instance.md)	 - Provides functionality for Argus instances

