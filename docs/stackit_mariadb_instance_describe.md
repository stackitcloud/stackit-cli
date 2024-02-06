## stackit mariadb instance describe

Shows details  of a MariaDB instance

### Synopsis

Shows details  of a MariaDB instance.

```
stackit mariadb instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a MariaDB instance with ID "xxx"
  $ stackit mariadb instance describe xxx

  Get details of a MariaDB instance with ID "xxx" in a table format
  $ stackit mariadb instance describe xxx --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit mariadb instance describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mariadb instance](./stackit_mariadb_instance.md)	 - Provides functionality for MariaDB instances

