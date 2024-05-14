## stackit mariadb instance delete

Deletes a MariaDB instance

### Synopsis

Deletes a MariaDB instance.

```
stackit mariadb instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a MariaDB instance with ID "xxx"
  $ stackit mariadb instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit mariadb instance delete"
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

* [stackit mariadb instance](./stackit_mariadb_instance.md)	 - Provides functionality for MariaDB instances

