## stackit mongodbflex user reset-password

Resets the password of a MongoDB Flex user

### Synopsis

Resets the password of a MongoDB Flex user.
sThe new password is visible after and cannot be retrieved later.

```
stackit mongodbflex user reset-password USER_ID [flags]
```

### Examples

```
  Reset the password of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit mongodbflex user reset-password xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex user reset-password"
      --instance-id string   ID of the instance
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit mongodbflex user](./stackit_mongodbflex_user.md)	 - Provides functionality for MongoDB Flex users

