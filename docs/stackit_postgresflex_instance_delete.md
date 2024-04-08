## stackit postgresflex instance delete

Deletes a PostgreSQL Flex instance

### Synopsis

Deletes a PostgreSQL Flex instance.

```
stackit postgresflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit postgresflex instance delete"
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

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

