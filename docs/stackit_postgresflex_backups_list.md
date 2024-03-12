## stackit postgresflex backups list

Lists all backups which are available for a specific PostgreSQL Flex instance

### Synopsis

Lists all backups which are available for a specific PostgreSQL Flex instance.

```
stackit postgresflex backups list [flags]
```

### Examples

```
  List all backups of instance with ID "xxx"
  $ stackit postgresflex backups list xxx

  List all backups of instance with ID "xxx" in JSON format
  $ stackit postgresflex backups list xxx --output-format json

  List up to 10 backups of instance with ID "xxx"
  $ stackit postgresflex backups list xxx --limit 10
```

### Options

```
  -h, --help        Help for "stackit postgresflex backups list"
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

* [stackit postgresflex backups](./stackit_postgresflex_backups.md)	 - Provides functionality for PostgreSQL Flex instance backups

