## stackit mariadb instance list

Lists all MariaDB instances

### Synopsis

Lists all MariaDB instances.

```
stackit mariadb instance list [flags]
```

### Examples

```
  List all MariaDB instances
  $ stackit mariadb instance list

  List all MariaDB instances in JSON format
  $ stackit mariadb instance list --output-format json

  List up to 10 MariaDB instances
  $ stackit mariadb instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit mariadb instance list"
      --limit int   Maximum number of entries to list
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

* [stackit mariadb instance](./stackit_mariadb_instance.md)	 - Provides functionality for MariaDB instances

