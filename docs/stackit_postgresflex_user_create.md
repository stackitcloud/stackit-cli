## stackit postgresflex user create

Create a PostgreSQL Flex user

### Synopsis

Create a PostgreSQL Flex user.
The password is only visible upon creation and cannot be retrieved later.
Alternatively, you can reset the password and access the new one by running:
  $ stackit postgresflex user reset-password --instance-id <INSTANCE_ID> --user-id <USER_ID>

```
stackit postgresflex user create [flags]
```

### Examples

```
  Create a PostgreSQL Flex user for instance with ID "xxx" and specify the username
  $ stackit postgresflex user create --instance-id xxx --username johndoe --roles read

  Create a PostgreSQL Flex user for instance with ID "xxx" with an automatically generated username
  $ stackit postgresflex user create --instance-id xxx --roles read
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user create"
      --instance-id string   ID of the instance
      --roles strings        Roles of the user, possible values are ["read" "readWrite"] (default [read])
      --username string      Username of the user. If not specified, a random username will be assigned
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

