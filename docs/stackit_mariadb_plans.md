## stackit mariadb plans

Lists all MariaDB service plans

### Synopsis

Lists all MariaDB service plans.

```
stackit mariadb plans [flags]
```

### Examples

```
  List all MariaDB service plans
  $ stackit mariadb plans

  List all MariaDB service plans in JSON format
  $ stackit mariadb plans --output-format json

  List up to 10 MariaDB service plans
  $ stackit mariadb plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit mariadb plans"
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

* [stackit mariadb](./stackit_mariadb.md)	 - Provides functionality for MariaDB

