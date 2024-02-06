## stackit mariadb credentials delete

Deletes credentials of a MariaDB instance

### Synopsis

Deletes credentials of a MariaDB instance.

```
stackit mariadb credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete credentials with ID "xxx" of MariaDB instance with ID "yyy"
  $ stackit mariadb credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit mariadb credentials delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mariadb credentials](./stackit_mariadb_credentials.md)	 - Provides functionality for MariaDB credentials

