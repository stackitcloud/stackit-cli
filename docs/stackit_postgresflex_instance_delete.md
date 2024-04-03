## stackit postgresflex instance delete

Deletes a PostgreSQL Flex instance

### Synopsis

Deletes a PostgreSQL Flex instance.
By default, instances will be kept in a delayed deleted state for 7 days before being permanently deleted.
Use the --force flag to force the deletion of a delayed deleted instance.

```
stackit postgresflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex instance delete xxx

  Force the deletion of a delayed deleted PostgreSQL Flex instance with ID "xxx"
  $ stackit postgresflex instance delete xxx --force
```

### Options

```
  -f, --force   Force deletion of a delayed deleted instance
  -h, --help    Help for "stackit postgresflex instance delete"
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

