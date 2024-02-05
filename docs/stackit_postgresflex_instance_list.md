## stackit postgresflex instance list

List all PostgreSQL Flex instances

### Synopsis

List all PostgreSQL Flex instances.

```
stackit postgresflex instance list [flags]
```

### Examples

```
  List all PostgreSQL Flex instances
  $ stackit postgresflex instance list

  List all PostgreSQL Flex instances in JSON format
  $ stackit postgresflex instance list --output-format json

  List up to 10 PostgreSQL Flex instances
  $ stackit postgresflex instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit postgresflex instance list"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

