## stackit mariadb credentials create

Creates credentials for a MariaDB instance

### Synopsis

Creates credentials (username and password) for a MariaDB instance.

```
stackit mariadb credentials create [flags]
```

### Examples

```
  Create credentials for a MariaDB instance
  $ stackit mariadb credentials create --instance-id xxx

  Create credentials for a MariaDB instance and show the password in the output
  $ stackit mariadb credentials create --instance-id xxx --show-password
```

### Options

```
  -h, --help                 Help for "stackit mariadb credentials create"
      --instance-id string   Instance ID
  -s, --show-password        Show password in output
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

* [stackit mariadb credentials](./stackit_mariadb_credentials.md)	 - Provides functionality for MariaDB credentials

