## stackit postgresflex user delete

Deletes a PostgreSQL Flex user

### Synopsis

Deletes a PostgreSQL Flex user by ID.
You can get the IDs of users for an instance by running:
  $ stackit postgresflex user list --instance-id <INSTANCE_ID>

```
stackit postgresflex user delete USER_ID [flags]
```

### Examples

```
  Delete a PostgreSQL Flex user with ID "xxx" for instance with ID "yyy"
  $ stackit postgresflex user delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user delete"
      --instance-id string   Instance ID
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

* [stackit postgresflex user](./stackit_postgresflex_user.md)	 - Provides functionality for PostgreSQL Flex users

