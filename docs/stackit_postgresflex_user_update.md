## stackit postgresflex user update

Updates a PostgreSQL Flex user

### Synopsis

Updates a PostgreSQL Flex user.

```
stackit postgresflex user update USER_ID [flags]
```

### Examples

```
  Update the roles of a PostgreSQL Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit postgresflex user update xxx --instance-id yyy --roles read
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user update"
      --instance-id string   ID of the instance
      --roles strings        Roles of the user, possible values are ["read" "readWrite"] (default [])
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

