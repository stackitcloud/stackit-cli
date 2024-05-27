## stackit sqlserverflex instance delete

Deletes an SQLServer Flex instance

### Synopsis

Deletes an SQLServer Flex instance.

```
stackit sqlserverflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete an SQLServer Flex instance with ID "xxx"
  $ stackit sqlserverflex instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit sqlserverflex instance delete"
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

* [stackit sqlserverflex instance](./stackit_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

