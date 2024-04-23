## stackit mariadb credentials describe

Shows details of credentials of a MariaDB instance

### Synopsis

Shows details of credentials of a MariaDB instance. The password will be shown in plain text in the output.

```
stackit mariadb credentials describe CREDENTIALS_ID [flags]
```

### Examples

```
  Get details of credentials with ID "xxx" from instance with ID "yyy"
  $ stackit mariadb credentials describe xxx --instance-id yyy

  Get details of credentials with ID "xxx" from instance with ID "yyy" in a table format
  $ stackit mariadb credentials describe xxx --instance-id yyy --output-format pretty
```

### Options

```
  -h, --help                 Help for "stackit mariadb credentials describe"
      --instance-id string   Instance ID
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

* [stackit mariadb credentials](./stackit_mariadb_credentials.md)	 - Provides functionality for MariaDB credentials

