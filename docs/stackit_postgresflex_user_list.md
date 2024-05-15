## stackit postgresflex user list

Lists all PostgreSQL Flex users of an instance

### Synopsis

Lists all PostgreSQL Flex users of an instance.

```
stackit postgresflex user list [flags]
```

### Examples

```
  List all PostgreSQL Flex users of instance with ID "xxx"
  $ stackit postgresflex user list --instance-id xxx

  List all PostgreSQL Flex users of instance with ID "xxx" in JSON format
  $ stackit postgresflex user list --instance-id xxx --output-format json

  List up to 10 PostgreSQL Flex users of instance with ID "xxx"
  $ stackit postgresflex user list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
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

* [stackit postgresflex user](./stackit_postgresflex_user.md)	 - Provides functionality for PostgreSQL Flex users

