## stackit argus instance describe

Shows details of an Argus instance

### Synopsis

Shows details of an Argus instance.

```
stackit argus instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of an Argus instance with ID "xxx"
  $ stackit argus instance describe xxx

  Get details of an Argus instance with ID "xxx" in a table format
  $ stackit argus instance describe xxx --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit argus instance describe"
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

