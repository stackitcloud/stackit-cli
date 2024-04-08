## stackit argus instance list

Lists all Argus instances

### Synopsis

Lists all Argus instances.

```
stackit argus instance list [flags]
```

### Examples

```
  List all Argus instances
  $ stackit argus instance list

  List all Argus instances in JSON format
  $ stackit argus instance list --output-format json

  List up to 10 Argus instances
  $ stackit argus instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit argus instance list"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus instance](./stackit_argus_instance.md)	 - Provides functionality for Argus instances

