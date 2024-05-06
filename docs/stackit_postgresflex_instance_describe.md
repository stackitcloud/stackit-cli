## stackit postgresflex instance describe

Shows details of a PostgreSQL Flex instance

### Synopsis

Shows details of a PostgreSQL Flex instance.

```
stackit postgresflex instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex instance describe xxx

  Get details of a PostgreSQL Flex instance with ID "xxx" in JSON format
  $ stackit postgresflex instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit postgresflex instance describe"
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

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

