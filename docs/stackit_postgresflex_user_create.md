## stackit postgresflex user create

Creates a PostgreSQL Flex user

### Synopsis

Creates a PostgreSQL Flex user.
The password is only visible upon creation and cannot be retrieved later.
Alternatively, you can reset the password and access the new one by running:
  $ stackit postgresflex user reset-password USER_ID --instance-id INSTANCE_ID

```
stackit postgresflex user create [flags]
```

### Examples

```
  Create a PostgreSQL Flex user for instance with ID "xxx"
  $ stackit postgresflex user create --instance-id xxx --username johndoe

  Create a PostgreSQL Flex user for instance with ID "xxx" and permission "createdb"
  $ stackit postgresflex user create --instance-id xxx --username johndoe --role createdb
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user create"
      --instance-id string   ID of the instance
      --role strings         Roles of the user, possible values are ["login" "createdb"] (default [login])
      --username string      Username of the user
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

* [stackit postgresflex user](./stackit_postgresflex_user.md)	 - Provides functionality for PostgreSQL Flex users

