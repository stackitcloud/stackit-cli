## stackit postgresflex user reset-password

Reset the password of a PostgreSQL Flex user

### Synopsis

Reset the password of a PostgreSQL Flex user. The new password is returned in the response.

```
stackit postgresflex user reset-password USER_ID [flags]
```

### Examples

```
  Reset the password of a PostgreSQL Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit postgresflex user reset-password xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user reset-password"
      --instance-id string   ID of the instance
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex user](./stackit_postgresflex_user.md)	 - Provides functionality for PostgreSQL Flex users

