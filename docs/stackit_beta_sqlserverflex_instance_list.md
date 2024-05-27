## stackit beta sqlserverflex instance list

Lists all SQLServer Flex instances

### Synopsis

Lists all SQLServer Flex instances.

```
stackit beta sqlserverflex instance list [flags]
```

### Examples

```
  List all SQLServer Flex instances
  $ stackit sqlserverflex instance list

  List all SQLServer Flex instances in JSON format
  $ stackit sqlserverflex instance list --output-format json

  List up to 10 SQLServer Flex instances
  $ stackit sqlserverflex instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta sqlserverflex instance list"
      --limit int   Maximum number of entries to list
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

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

