## stackit beta sqlserverflex user reset-password

Resets the password of a SQLServer Flex user

### Synopsis

Resets the password of a SQLServer Flex user.
sThe new password is visible after and cannot be retrieved later.

```
stackit beta sqlserverflex user reset-password USER_ID [flags]
```

### Examples

```
  Reset the password of a SQLServer Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit sqlserverflex user reset-password xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit beta sqlserverflex user reset-password"
      --instance-id string   ID of the instance
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

* [stackit beta sqlserverflex user](./stackit_beta_sqlserverflex_user.md)	 - Provides functionality for SQLServer Flex users

